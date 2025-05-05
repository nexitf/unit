package plan

import (
	"context"
	"time"
)

type TimerRoutine struct {
	d  time.Duration
	fn func(time.Time)
}

// NewTimerRoutine returns a new ticker routine.
//
// The parameter `d` specifies the interval at which the timer runs.
// The parameter `fn` is a callback function.
func NewTimerRoutine(d time.Duration, fn func(time.Time)) (r *TimerRoutine) {
	return &TimerRoutine{d: d, fn: fn}
}

// Name implements core.Routine.
func (r *TimerRoutine) Name() string {
	return "NexITF timer plan"
}

// Run implements core.Routine.
func (r *TimerRoutine) Run(ctx context.Context) (err error) {
	tm := time.NewTimer(r.d)
	defer func() {
		tm.Stop()
	}()
	for {
		select {
		// Cancel context
		case <-ctx.Done():
			return ctx.Err()
		// Invoke
		case tm := <-tm.C:
			r.fn(tm)
			return
		}
	}
}

// Stop implements core.Routine.
func (r *TimerRoutine) Stop(ctx context.Context) (err error) {
	return
}
