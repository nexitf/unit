package unit

import (
	"time"

	"github.com/nexitf/unit/internal/core"
)

type RunOption = core.RunOption

// WithDelayReady
func WithDelayReady(d time.Duration) RunOption {
	return core.WithDelayReady(d)
}
