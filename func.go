package unit

import (
	"context"
	"fmt"
	"time"
)

type FuncRoutine struct {
	id int64
	fn func(ctx context.Context) (err error)
}

// NewFuncRoutine
func NewFuncRoutine(fn func(ctx context.Context) (err error)) (r *FuncRoutine) {
	return &FuncRoutine{id: time.Now().Unix(), fn: fn}
}

// Name implements routine.Routine.
func (r *FuncRoutine) Name() string {
	return fmt.Sprintf("func-%d", r.id)
}

// Run implements routine.Routine.
func (r *FuncRoutine) Run(ctx context.Context) (err error) {
	return r.fn(ctx)
}
