package unit

import (
	"context"
	"encoding/json"

	"github.com/nexitf/unit/internal/core"
)

const (
	Version = core.Version
)

type Routine = core.Routine
type ReadyChecker = core.ReadyChecker

// Setup
func Setup(routine Routine) {
	core.Setup(routine)
}

// Inspect returns all infomation of the current unit in JSON format.
func Inspect() string {
	rt := core.Inspect()
	out, err := json.MarshalIndent(rt, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(out)
}

// Run
func Run(ctx context.Context) (err error) {
	return core.Run(ctx)
}
