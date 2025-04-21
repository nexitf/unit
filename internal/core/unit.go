package core

import (
	"context"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/nexitf/unit/internal/core/plugin"
	"github.com/nexitf/unit/internal/core/utils"
)

var (
	// Unit instance
	u Unit
)

type Unit struct {
	routines   []Routine                // Routines
	externals  map[string]Connector     // Dependent external resources
	plugins    map[string]plugin.Plugin // Plugins
	unitExited chan struct{}
	procExited chan os.Signal
	delayReady time.Duration
	cancelCtx  context.CancelFunc
	running    int32 // Already running
	fatal      error
	errored    int32 // Error found
}

func init() {
	u.routines = make([]Routine, 0)
	u.externals = make(map[string]Connector)
	u.unitExited = make(chan struct{})
	u.procExited = make(chan os.Signal, 1)
}

// Init
func (u *Unit) Init(ctx context.Context) (err error) {
	if u.delayReady <= 0 {
		u.delayReady = 50 * time.Millisecond
	}
	// Load all plugins
	u.plugins = plugin.LoadPlugins()
	return
}

// Setup
func (u *Unit) Setup(routine Routine) {
	u.routines = append(u.routines, &Launcher{routine: routine})
}

// Setup
func Setup(routine Routine) {
	if u.IsRunning() {
		panic(ErrAlreadyRunning)
	}
	u.Setup(routine)
}

// WaitForExit
func (u *Unit) WaitForExit(ctx context.Context) (err error) {
	var (
		ctxDone = ctx.Done()
		timer   = utils.NewStoppedTimer()
	)

	defer timer.Stop()

	// Register the signals to be monitored: interrupt (Ctrl+C) and termination
	signal.Notify(u.procExited, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		// Cancel context
		case <-ctxDone:
			u.Panic(ctx.Err())
			timer.Reset(3 * time.Second)
			ctxDone = nil
		// Delay exited
		case <-timer.C:
			return u.Fatal()
		// Process exited
		case <-u.procExited:
			u.Panic(ErrProcessExited)
			u.cancelCtx()
			timer.Reset(3 * time.Second)
		// Unit exited
		case <-u.unitExited:
			return u.Fatal()
		}
	}
}

type RunOption func(*Unit)

// WithDelayReady
func WithDelayReady(d time.Duration) RunOption {
	return func(u *Unit) {
		if d > 0 {
			u.delayReady = d
		}
	}
}

// Run
func (u *Unit) Run(ctx context.Context, opts ...RunOption) (err error) {
	// No routine
	if len(u.routines) <= 0 {
		return
	}

	// Set options
	for _, setOpt := range opts {
		setOpt(u)
	}

	atomic.StoreInt32(&u.running, 1)
	defer atomic.StoreInt32(&u.running, 0)

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

	// Start routines
	u.StartRoutines(ctx)

	readyFn := func() {
		if !u.Errored() {
			if err := u.ReadyRoutines(ctx); err != nil {
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

// Panic
func (u *Unit) Panic(err error) {
	if u.fatal != nil {
		return
	}
	u.fatal = err
	atomic.StoreInt32(&u.errored, 1)
}

// Fatal returns the fatal error
// encountered during the unit's lifecycle.
func (u *Unit) Fatal() (err error) {
	return u.fatal
}

// Run
func Run(ctx context.Context, opts ...RunOption) (err error) {
	if u.IsRunning() {
		panic(ErrAlreadyRunning)
	}
	return u.Run(ctx, opts...)
}
