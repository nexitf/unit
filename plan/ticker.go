package plan

import (
	"context"
	"time"
)

type TickerRoutine struct {
	d  time.Duration
	fn func(time.Time)
}

// NewTickerRoutine returns a new ticker routine.
//
// The parameter `d` specifies the interval at which the ticker runs.
// The parameter `fn` is a callback function that is invoked on each tick.
func NewTickerRoutine(d time.Duration, fn func(time.Time)) (r *TickerRoutine) {
	return &TickerRoutine{d: d, fn: fn}
}

// Name implements core.Routine.
func (r *TickerRoutine) Name() string {
	return "NexITF ticker plan"
}

// Run implements core.Routine.
func (r *TickerRoutine) Run(ctx context.Context) (err error) {
	tick := time.NewTicker(r.d)
	defer func() {
		tick.Stop()
	}()
	for {
		select {
		// Cancel context
		case <-ctx.Done():
			return ctx.Err()
		// Invoke
		case tm := <-tick.C:
			r.fn(tm)
		}
	}
}

// Stop implements core.Routine.
func (r *TickerRoutine) Stop(ctx context.Context) (err error) {
	return
}
