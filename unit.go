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
type RoutineChecker = routine.RoutineChecker

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
	core.UsePlugin(plugins...)
}

// Run
func Run(ctx context.Context, opts ...RunOption) (err error) {
	return core.Run(ctx, opts...)
}
