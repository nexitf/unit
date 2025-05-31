package unit

import (
	"time"

	"github.com/nexitf/unit/internal/core"
)

type RunOption = core.RunOption

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
	return core.WithDelayReady(d)
}

// WithBindTimeout creates a new RunOption that sets the timeout duration for
// binding external resources during the runtime.
// If the provided duration 'd' is greater than zero, it will be assigned to the 'bindTimeout' field of the Unit.
// Otherwise, the 'bindTimeout' field remains unchanged.
//
// Parameters:
//   - d: The timeout duration for binding external resources.
//
// Returns:
//   - RunOption: A functional option that can be used to configure a Unit.
func WithBindTimeout(d time.Duration) RunOption {
	return core.WithBindTimeout(d)
}
