package unit

import (
	"context"

	"github.com/nexitf/unit/internal/core"
	"github.com/nexitf/unit/internal/core/plugin"
	"github.com/nexitf/unit/internal/core/routine"
)

const (
	Version = core.Version
)

type Plugin = plugin.Plugin
type Routine = routine.Routine
type RoutineStopped = routine.RoutineStopped
type RoutineChecker = routine.RoutineChecker

// RoutineGroup implements the Routine interface and is used to set routines in batches.
// RoutineGroup itself will not be executed as a routine.
type RoutineGroup = core.RoutineGroup

// NewRoutineGroup creates a new instance of RoutineGroup.
// It initializes an empty slice of routines and returns a pointer to the newly created RoutineGroup.
// This function serves as a convenient way to instantiate a RoutineGroup.
func NewRoutineGroup() (rg *RoutineGroup) {
	return core.NewRoutineGroup()
}

// Init sets a init function to be executed when the unit inits.
func Init(fn func(ctx context.Context)) {
	core.Init(fn)
}

// Defer sets a callback function to be executed when the unit exits.
func Defer(fn func()) {
	core.Defer(fn)
}

// Use adds plugin to the Unit instance.
// If the unit is already running, it logs an error and panics.
func Use(plugins ...Plugin) {
	core.Use(plugins...)
}

// Run starts and runs a specified routine.
// This function calls the Run method of the internal core module to actually execute the routine.
//
// Parameters:
//   - ctx is a context object used to control the lifecycle of the routine, which can be used for cancellation or timeout control.
//   - routine is the routine instance to be run, which must implement the routine.Routine interface.
//   - opts are optional running options used to configure the running behavior of the routine.
//
// Returns
//   - The return value err represents any errors that may occur during the running process.
func Run(ctx context.Context, routine Routine, opts ...RunOption) (err error) {
	return core.Run(ctx, routine, opts...)
}
