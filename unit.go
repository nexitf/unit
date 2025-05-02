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

// Setup sets a routine implementation that will be launched and
// run as an instance during program execution.
func Setup(routine Routine) {
	core.Setup(routine)
}

// Init sets a init function to be executed when the unit inits.
func Init(fn func(ctx context.Context)) {
	core.Init(fn)
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

// WithDelayReady
func WithDelayReady(d time.Duration) RunOption {
	return core.WithDelayReady(d)
}

// Run
func Run(ctx context.Context, opts ...RunOption) (err error) {
	return core.Run(ctx, opts...)
}
