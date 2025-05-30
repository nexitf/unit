package routine

import (
	"context"

	"github.com/nexitf/unit/internal/errors"
)

var (
	ErrRoutineException = errors.New("routine exception")
)

type Routine interface {
	// Name returns the name of the routine.
	Name() string

	// Run Usually, Run is the entry point of a routine.
	// A running routine is expected to block inside Run.
	Run(ctx context.Context) (err error)
}

type RoutineStopped interface {
	// Stop method will be called once when the routine exits.
	Stop(ctx context.Context) (err error)
}

type RoutineChecker interface {
	Ready(ctx context.Context) (err error)
}
