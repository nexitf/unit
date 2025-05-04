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
	inits      []func(context.Context)
	defers     []func()
}

func init() {
	u.routines = make([]Routine, 0)
	u.externals = make(map[string]Connector)
	u.unitExited = make(chan struct{})
	u.procExited = make(chan os.Signal, 1)
}

// Init
func (u *Unit) init(ctx context.Context) (err error) {
	if u.delayReady <= 0 {
		u.delayReady = 50 * time.Millisecond
	}
	// Load all plugins
	u.plugins = plugin.LoadPlugins()
	return
}

// Setup sets a routine implementation that will be launched and
// run as an instance during program execution.
func (u *Unit) Setup(routine Routine) {
	u.routines = append(u.routines, &Launcher{routine: routine})
}

// Setup sets a routine implementation that will be launched and
// run as an instance during program execution.
func Setup(routine Routine) {
	if u.IsRunning() {
		logkit.PanicWrap(ErrAlreadyRunning, "unit already running")
		panic(ErrAlreadyRunning)
	}
	u.Setup(routine)
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

// WaitForExit
func (u *Unit) WaitForExit(ctx context.Context) (err error) {
	var (
		ctxDone = ctx.Done()
		timer   = utils.NewNeverTriggeringTimer()
	)

	defer timer.Stop()

	logkit.Info("wait for exit")

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
		case signal := <-u.procExited:
			logkit.Warn("process exited", logkit.Field("signal", signal.String()))
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
		logkit.Warn("no routine found")
		return
	}

	// Set options
	for _, setOpt := range opts {
		setOpt(u)
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

	routineCtx, cancelRoutineCtx := context.WithCancel(ctx)
	defer cancelRoutineCtx()

	u.cancelCtx = cancelRoutineCtx

	// Start routines
	u.StartRoutines(routineCtx)

	logkit.Info("all routines started successfully")

	readyFn := func() {
		if !u.Errored() {
			if err := u.ReadyRoutines(ctx); err != nil {
				logkit.ErrorWrap(err, "routines not ready")
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
		logkit.PanicWrap(ErrAlreadyRunning, "unit already running")
		panic(ErrAlreadyRunning)
	}
	return u.Run(ctx, opts...)
}
