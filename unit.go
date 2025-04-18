package unit

import (
	"context"

	"github.com/nexitf/unit/internal/core"
)

type Task = core.Task

// Name
func Name() string {
	return core.Name()
}

// Level
func Level() int {
	return core.Level()
}

// Setup
func Setup(task Task) {
	core.Setup(task)
}

// Run
func Run(ctx context.Context) (err error) {
	return core.Run(ctx)
}
