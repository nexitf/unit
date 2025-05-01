package core

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nexitf/logkit"
)

type Routine interface {
	// Name returns the name of the routine.
	Name() string

	// Run Usually, Run is the entry point of a routine.
	// A running routine is expected to block inside Run.
	Run(ctx context.Context) (err error)

	// Stop method will be called once when the routine exits.
	Stop(ctx context.Context) (err error)
}

type Launcher struct {
	routine  Routine
	err      error
	runTime  time.Time
	stopTime time.Time
}

// Name implements Routine.
func (l *Launcher) Name() string {
	return l.routine.Name()
}

// Run implements Routine.
func (l *Launcher) Run(ctx context.Context) (err error) {
	defer func() {
		if l.err == nil {
			l.err = err
		}
		l.stopTime = time.Now()
	}()
	l.runTime = time.Now()
	return l.routine.Run(ctx)
}

// Stop implements Routine.
func (l *Launcher) Stop(ctx context.Context) (err error) {
	l.stopTime = time.Now()
	return l.routine.Stop(ctx)
}

// Ready implements RoutineChecker.
func (l *Launcher) Ready(ctx context.Context) (err error) {
	rc, ok := l.routine.(RoutineChecker)
	if !ok {
		return nil
	}
	defer func() {
		if l.err == nil {
			l.err = err
		}
	}()
	return rc.Ready(ctx)
}

type RoutineChecker interface {
	Ready(ctx context.Context) (err error)
}

// StartRoutines
func (u *Unit) StartRoutines(ctx context.Context) {
	var (
		first   int32
		started sync.WaitGroup
		finish  sync.WaitGroup
	)

	// Launch tasks
	for _, r := range u.routines {
		started.Add(1)
		finish.Add(1)
		go func(r Routine) {
			started.Done()
			// Run routine
			if err := r.Run(ctx); err == nil {
				r.Stop(ctx)
			} else {
				logkit.ErrorWrap(err, "routine run failed", logkit.Field("routine", r.Name()))
				u.Panic(err)
			}
			finish.Done()
			atomic.StoreInt32(&first, 1)
		}(r)
	}

	// Wait for all tasks be started
	started.Wait()

	go func() {
		// Will be blocked here
		for {
			if atomic.LoadInt32(&first) == 1 {
				break
			} else {
				runtime.Gosched()
			}
		}

		if u.Errored() {
			// Exit all routines
			u.cancelCtx()
		}

		// Wait for all tasks to exit
		finish.Wait()
		close(u.unitExited)
	}()
}

// ReadyRoutines
func (u *Unit) ReadyRoutines(ctx context.Context) (err error) {
	for _, r := range u.routines {
		rc, ok := r.(RoutineChecker)
		if !ok {
			continue
		}
		if err = rc.Ready(ctx); err != nil {
			logkit.ErrorWrap(err, "routine readiness check failed", logkit.Field("routine", r.Name()))
			u.Panic(err)
			return
		}
	}
	return
}
