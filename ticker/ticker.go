package ticker

import (
	"context"
	"time"
)

type Routine struct {
	d  time.Duration
	fn func(time.Time) bool
}

func NewRoutine(d time.Duration, fn func(time.Time) bool) (r *Routine) {
	r = &Routine{d: d, fn: fn}
	return
}

// Name implements core.Routine.
func (r *Routine) Name() string {
	return "NexITF ticker"
}

// Run implements core.Routine.
func (r *Routine) Run(ctx context.Context) (err error) {
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
func (r *Routine) Stop(ctx context.Context) (err error) {
	return
}
