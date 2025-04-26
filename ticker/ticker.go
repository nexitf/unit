package ticker

import (
	"context"
	"time"
)

type TickerRoutine struct {
	d  time.Duration
	fn func(time.Time) bool
}

// NewTickerRoutine returns a new ticker routine.
//
// The parameter `d` specifies the interval at which the ticker runs.
// The parameter `fn` is a callback function that is invoked on each tick.
// If `fn` returns true, the ticker continues running; otherwise, it stops.
func NewTickerRoutine(d time.Duration, fn func(time.Time) bool) (r *TickerRoutine) {
	return &TickerRoutine{d: d, fn: fn}
}

// Name implements core.Routine.
func (r *TickerRoutine) Name() string {
	return "NexITF ticker"
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
		case tm := <-tick.C:
			if !r.fn(tm) {
				// Exited
				return
			}
		}
	}
}

// Stop implements core.Routine.
func (r *TickerRoutine) Stop(ctx context.Context) (err error) {
	return
}
