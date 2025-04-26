package unit

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nexitf/unit/internal/core"
)

const (
	Version = core.Version
)

type RunOption = core.RunOption
type Routine = core.Routine
type RoutineChecker = core.RoutineChecker

// WithDelayReady
func WithDelayReady(d time.Duration) RunOption {
	return core.WithDelayReady(d)
}

// Setup
func Setup(routine Routine) {
	core.Setup(routine)
}

// Defer sets a callback function to be executed when the unit exits.
func Defer(fn func()) {
	core.Defer(fn)
}

type Runtime = core.Runtime

// Inspect returns all infomation of the current unit.
func Inspect() (rt Runtime) {
	return core.Inspect()
}

// InspectJSON returns all infomation of the current unit in JSON format.
func InspectJSON() string {
	rt := core.Inspect()
	out, err := json.MarshalIndent(rt, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(out)
}

// Run
func Run(ctx context.Context, opts ...RunOption) (err error) {
	return core.Run(ctx, opts...)
}
