package core

import (
	"context"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/nexitf/logkit"
	"github.com/nexitf/unit/internal/core/plugin"
	"github.com/nexitf/unit/internal/core/routine"
	"github.com/nexitf/unit/internal/core/utils"
)

var (
	// Unit instance
	u Unit
)

type Unit struct {
	runCtx    context.Context
	scheduler *routine.Scheduler
	externals []*Connector             // Dependent external resources
	exloadeds map[string]struct{}      // The external resources that have been loaded
	plugins   map[string]plugin.Plugin // Plugins
	// Status
	running int32 // Already running
	fatal   error
	inits   []func(context.Context)
	defers  []func()
	// Options
	delayReady  time.Duration
	bindTimeout time.Duration
}

func init() {
	u.scheduler = routine.NewScheduler()
	u.exloadeds = make(map[string]struct{})
	u.externals = make([]*Connector, 0)
	u.plugins = make(map[string]plugin.Plugin)
}

// Init
func (u *Unit) init(_ context.Context) (err error) {
	if u.delayReady <= 0 {
		u.delayReady = 50 * time.Millisecond
	}
	if u.bindTimeout <= 0 {
		u.bindTimeout = 50 * time.Millisecond
	}
	return
}

// Init sets a init function to be executed when the unit inits.
func (u *Unit) Init(fn func(context.Context)) {
	if u.IsRunning() {
		logkit.PanicWrap(ErrAlreadyRunning, "unit already running")
		panic(ErrAlreadyRunning)
	}
	u.inits = append(u.inits, fn)
}

// Init sets a init function to be executed when the unit inits.
func Init(fn func(context.Context)) {
	u.Init(fn)
}

// Defer sets a callback function to be executed when the unit exits.
func (u *Unit) Defer(fn func()) {
	u.defers = append(u.defers, fn)
}

// Defer sets a callback function to be executed when the unit exits.
func Defer(fn func()) {
	u.Defer(fn)
}

// Use adds plugin to the Unit instance.
// If the unit is already running, it logs an error and panics.
func (u *Unit) Use(plugins ...plugin.Plugin) {
	if u.IsRunning() {
		logkit.PanicWrap(ErrAlreadyRunning, "unit already running")
		panic(ErrAlreadyRunning)
	}
	for _, plugin := range plugins {
		pluginID, err := utils.RandomString(32)
		if err != nil {
			logkit.PanicWrap(err, "plugin ID generation failed")
			panic(err)
		}
		u.plugins[pluginID] = plugin
		logkit.Info("enable plugin", logkit.Field("plugin", plugin.Name()))
	}
}

// Use adds plugin to the Unit instance.
// If the unit is already running, it logs an error and panics.
func Use(plugins ...plugin.Plugin) {
	u.Use(plugins...)
}

// WaitForExit
func (u *Unit) WaitForExit(ctx context.Context) (err error) {
	var (
		ctxDone  = ctx.Done()
		procDone = make(chan os.Signal, 1)
		timer    = utils.NewNeverTriggeringTimer()
	)

	defer timer.Stop()

	logkit.Info("wait for exit")

	// Used to interrupt the execution of the routines
	interrupt := func(err error) {
		u.fatal = err
		u.scheduler.Interrupt(err)
		timer.Reset(3 * time.Second)
	}

	// Register the signals to be monitored: interrupt (Ctrl+C) and termination
	signal.Notify(procDone, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		// Cancel context
		case <-ctxDone:
			interrupt(ctx.Err())
			ctxDone = nil
		// Delay exited
		case <-timer.C:
			return u.fatal
		// Process exited
		case signal := <-procDone:
			interrupt(ErrProcessTerminated)
			logkit.Warn("process exited", logkit.Field("signal", signal.String()))
		// Unit exited
		case <-u.scheduler.Done():
			u.fatal = u.scheduler.Err()
			return u.fatal
		}
	}
}

// Run
func (u *Unit) Run(ctx context.Context, routine routine.Routine, opts ...RunOption) (err error) {
	u.runCtx = ctx

	// Set options
	for _, setOpt := range opts {
		setOpt(u)
	}

	// Load routine
	u.scheduler.Add(routine)

	// No routine
	if u.scheduler.Num() <= 0 {
		logkit.Warn("no routine found")
		return
	}

	atomic.StoreInt32(&u.running, 1)
	defer atomic.StoreInt32(&u.running, 0)

	logkit.Info("unit running")
	defer func() {
		if err != nil {
			logkit.ErrorWrap(err, "unit run failed")
		}
		logkit.Info("unit exited")
	}()

	// Init unit
	if err = u.init(ctx); err != nil {
		logkit.ErrorWrap(err, "no routine found")
		return
	}

	// Run plugin
	if err = u.RunPlugins(ctx); err != nil {
		logkit.ErrorWrap(err, "run plugins failed")
		return
	}
	defer u.StopPlugins(ctx)

	// Load dependent external resources
	if err = u.LoadExternals(ctx); err != nil {
		logkit.ErrorWrap(err, "load external resources failed")
		return
	}

	// All init callbacks are invoked once resource loading is complete.
	for _, fn := range u.inits {
		fn(ctx)
	}

	// Start routines
	u.scheduler.Start(ctx)
	logkit.Info("all routines started successfully")

	readyFn := func() {
		u.scheduler.Ready(ctx)
	}
	if u.delayReady > 0 {
		// Wait for a short period to verify
		// that all tasks are ready
		time.AfterFunc(u.delayReady, readyFn)
	} else {
		readyFn()
	}

	defer func() {
		// During runtime, some resources may need to be released
		// before the program fully exits.
		for _, fn := range u.defers {
			fn()
		}
	}()

	// Will be blocked until exit
	return u.WaitForExit(ctx)
}

// Run
func Run(ctx context.Context, routine routine.Routine, opts ...RunOption) (err error) {
	if u.IsRunning() {
		logkit.PanicWrap(ErrAlreadyRunning, "unit already running")
		panic(ErrAlreadyRunning)
	}
	return u.Run(ctx, routine, opts...)
}
