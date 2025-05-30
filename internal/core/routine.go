package core

import (
	"context"

	"github.com/nexitf/unit/internal/core/routine"
)

// RoutineGroup implements the Routine interface and is used to set routines in batches.
// RoutineGroup itself will not be executed as a routine.
type RoutineGroup struct {
	routines []routine.Routine
}

// NewRoutineGroup creates a new instance of RoutineGroup.
// It initializes an empty slice of routines and returns a pointer to the newly created RoutineGroup.
// This function serves as a convenient way to instantiate a RoutineGroup.
func NewRoutineGroup() *RoutineGroup {
	return &RoutineGroup{}
}

// Name returns the name of the RoutineGroup.
// This method implements the corresponding method in the routine.Routine interface,
// providing a standardized way to identify the RoutineGroup instance.
// It always returns the fixed string "NEXITF routine group".
//
// Returns:
// - string: The name of the RoutineGroup.
func (rg *RoutineGroup) Name() string {
	return "NEXITF routine group"
}

// Add appends a new routine to the RoutineGroup.
// It takes a routine that implements the routine.Routine interface as a parameter.
// The provided routine is then added to the internal slice of routines within the RoutineGroup.
// This method modifies the state of the RoutineGroup by expanding its list of managed routines.
func (rg *RoutineGroup) Add(routine routine.Routine) {
	rg.routines = append(rg.routines, routine)
}

// Run will never be executed.
func (rg *RoutineGroup) Run(ctx context.Context) (err error) {
	return
}
