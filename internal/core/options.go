package core

import (
	"time"
)

type RunOption func(*Unit)

// WithDelayReady creates a new RunOption that sets the delay duration after all routines already started.
// If the provided duration 'd' is greater than zero, it will be assigned to the 'delayReady' field of the Unit.
// Otherwise, the 'delayReady' field remains unchanged.
//
// Parameters:
// - d: The duration to delay the Unit's readiness.
//
// Returns:
// - RunOption: A functional option that can be used to configure a Unit.
func WithDelayReady(d time.Duration) RunOption {
	return func(u *Unit) {
		if d > 0 {
			u.delayReady = d
		}
	}
}
