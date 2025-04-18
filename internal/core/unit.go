package core

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/nexitf/unit/internal/core/plugin"
)

var (
	// Unit instance
	u unit
)

type unit struct {
	routines  []Routine                // Routines
	externals map[string]Connector     // Dependent external resources
	plugins   map[string]plugin.Plugin // Plugins

	noReady    int32
	delayReady time.Duration
	cancelCtx  context.CancelFunc
	unitExited chan error
	err        error
	finish     sync.WaitGroup
	first      sync.WaitGroup
}

func init() {
	u.routines = make([]Routine, 0)
	u.externals = make(map[string]Connector)
}

// Init
func (u *unit) Init(ctx context.Context) (err error) {
	u.unitExited = make(chan error, 1)
	// Load all plugins
	u.plugins = plugin.LoadPlugins()
	return
}

// Setup
func (u *unit) Setup(routine Routine) {
	u.routines = append(u.routines, routine)
	// Init delay ready checker
	_, ok := routine.(ReadyChecker)
	if ok && u.delayReady <= 0 {
		u.delayReady = time.Second
	}
}

// Setup
func Setup(routine Routine) {
	u.Setup(routine)
}

// WaitForExit
func (u *unit) WaitForExit(ctx context.Context) (err error) {

	var (
		procExited = make(chan os.Signal, 1)
	)

	// Register the signals to be monitored: interrupt (Ctrl+C) and termination
	signal.Notify(procExited, syscall.SIGINT, syscall.SIGTERM)

	select {
	// Cancel context
	case <-ctx.Done():
		if u.err != nil {
			return u.err
		}
		return ctx.Err()
	// Process exited
	case <-procExited:
		return ErrProcessExited
	// Unit exited
	case err = <-u.unitExited:
		return err
	}
}

// Run
func (u *unit) Run(ctx context.Context) (err error) {
	// No routine
	if len(u.routines) <= 0 {
		return
	}

	// Init unit
	if err = u.Init(ctx); err != nil {
		return
	}

	// Run plugin
	if err = u.RunPlugins(ctx); err != nil {
		return
	}
	defer u.StopPlugins(ctx)

	// Load dependent external resources
	if err = u.LoadExternal(ctx); err != nil {
		return
	}

	ctx, u.cancelCtx = context.WithCancel(ctx)
	defer u.cancelCtx()

	// Run routines
	u.StartRoutines(ctx)

	readyFn := func() {
		if atomic.LoadInt32(&u.noReady) == 0 {
			if u.err = u.ReadyRoutines(ctx); u.err != nil {
				u.cancelCtx()
			}
		}
	}
	if u.delayReady > 0 {
		// Wait for a short period to verify
		// that all tasks are ready
		time.AfterFunc(u.delayReady, readyFn)
	} else {
		readyFn()
	}

	// Will be blocked until exit
	return u.WaitForExit(ctx)
}

// Run
func Run(ctx context.Context) (err error) {
	return u.Run(ctx)
}
