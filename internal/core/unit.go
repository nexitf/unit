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
	tasks     []Task                   // Tasks
	externals map[string]Connector     // Dependent external resources
	plugins   map[string]plugin.Plugin // Plugins

	noReady    int32
	delayReady time.Duration
	cancelTask context.CancelFunc
	taskExited chan error
	err        error
	finish     sync.WaitGroup
	first      sync.WaitGroup
}

func init() {
	u.tasks = make([]Task, 0)
	u.externals = make(map[string]Connector)
}

// Init
func (u *unit) Init(ctx context.Context) (err error) {
	if u.delayReady <= 0 {
		u.delayReady = time.Second
	}
	u.taskExited = make(chan error, 1)
	// Load all plugins
	u.plugins = plugin.LoadPlugins()
	return
}

// Setup
func (u *unit) Setup(task Task) {
	u.tasks = append(u.tasks, task)
}

// Setup
func Setup(task Task) {
	u.Setup(task)
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
		return ctx.Err()
	// Process exited
	case <-procExited:
		return ErrProcessExited
	// All tasks exited
	case err = <-u.taskExited:
		return err
	}
}

// Run
func (u *unit) Run(ctx context.Context) (err error) {
	// No task
	if len(u.tasks) <= 0 {
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

	ctx, u.cancelTask = context.WithCancel(ctx)
	defer u.cancelTask()

	// Run tasks
	u.StartTasks(ctx)

	// Wait for a short period to verify
	// that all tasks are ready
	time.AfterFunc(u.delayReady, func() {
		if atomic.LoadInt32(&u.noReady) == 0 {
			u.ReadyTasks(ctx)
		}
	})

	// Will be blocked until exit
	return u.WaitForExit(ctx)
}

// Run
func Run(ctx context.Context) (err error) {
	return u.Run(ctx)
}
