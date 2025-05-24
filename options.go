package unit

import (
	"time"

	"github.com/nexitf/unit/internal/core"
)

type RunOption = core.RunOption

// WithRoutine
func WithRoutine(routines ...Routine) RunOption {
	return core.WithRoutine(routines...)
}

// WithDelayReady
func WithDelayReady(d time.Duration) RunOption {
	return core.WithDelayReady(d)
}
