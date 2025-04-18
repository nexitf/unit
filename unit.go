package unit

import (
	"context"

	"github.com/nexitf/unit/internal/core"
)

type Routine = core.Routine
type ReadyChecker = core.ReadyChecker

// Setup
func Setup(routine Routine) {
	core.Setup(routine)
}

// Run
func Run(ctx context.Context) (err error) {
	return core.Run(ctx)
}
