package core

import (
	"context"
	"sync"
	"sync/atomic"
)

type Routine interface {
	Name() string
	Run(ctx context.Context) (err error)
	Stop(ctx context.Context) (err error)
}

type ReadyChecker interface {
	Ready(ctx context.Context) (err error)
}

// StartRoutines
func (u *unit) StartRoutines(ctx context.Context) {
	var (
		started sync.WaitGroup
	)

	u.first.Add(1)
	// Launch tasks
	for _, r := range u.routines {
		started.Add(1)
		u.finish.Add(1)
		go func(r Routine) {
			started.Done()
			u.err = r.Run(ctx)
			if u.err == nil {
				r.Stop(ctx)
			}
			u.first.Done()
			u.finish.Done()
		}(r)
	}

	// Wait for all tasks be started
	started.Wait()

	go func() {
		// Will be blocked here
		u.first.Wait()
		if u.err != nil {
			// Exit all tasks and stop sending readiness notifications
			// when the first runtime error is encountered
			atomic.StoreInt32(&u.noReady, 1)
			u.cancelCtx()
		}

		// Wait for all tasks to exit
		u.finish.Wait()
		u.unitExited <- u.err
	}()
}

// ReadyRoutines
func (u *unit) ReadyRoutines(ctx context.Context) (err error) {
	for _, r := range u.routines {
		health, ok := r.(ReadyChecker)
		if !ok {
			continue
		}
		if err = health.Ready(ctx); err != nil {
			return
		}
	}
	return
}
