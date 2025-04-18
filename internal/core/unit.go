package core

import (
	"context"
	"sync"

	"github.com/nexitf/unit/internal/core/plugin"
)

var (
	// Unit instance
	u unit
)

type Task interface {
	Name() string
	Run(ctx context.Context) (err error)
	Stop(ctx context.Context) (err error)
}

type unit struct {
	name      string                   // Name of unit
	level     int                      // Level of unit
	tasks     []Task                   // Tasks
	externals map[string]Connector     // Dependent external resources
	plugins   map[string]plugin.Plugin // Plugins
}

func init() {
	u.tasks = make([]Task, 0)
	u.externals = make(map[string]Connector)
}

// Init
func (u *unit) Init(ctx context.Context) (err error) {
	u.name = "unit"
	u.level = 0
	// Load all plugin plugins
	u.plugins = plugin.LoadPlugins()
	return
}

// Setup
func (u *unit) Setup(task Task) {
	u.tasks = append(u.tasks, task)
}

// RunTask
func (u *unit) RunTask(ctx context.Context) (
	first, finish *sync.WaitGroup, getLastErrFunc func() error) {
	var err error
	// Get last error
	getLastErrFunc = func() error {
		return err
	}

	first = &sync.WaitGroup{}
	finish = &sync.WaitGroup{}

	first.Add(1)
	// Launch tasks
	for _, task := range u.tasks {
		finish.Add(1)
		go func(task Task) {
			runErr := task.Run(ctx)
			if runErr != nil {
				err = runErr
			} else {
				task.Stop(ctx)
			}
			first.Done()
			finish.Done()
		}(task)
	}

	return
}

// Name returns unit name
func Name() string {
	return u.name
}

// Level returns unit level
func Level() int {
	return u.level
}

// Setup
func Setup(task Task) {
	u.Setup(task)
}

// Run
func Run(ctx context.Context) (err error) {
	// No task
	if len(u.tasks) <= 0 {
		return
	}

	// Init unit
	if err = u.Init(ctx); err != nil {
		return
	}

	// // Run plugin
	// if err = u.RunPlugin(ctx); err != nil {
	// 	return
	// }
	// defer u.StopPlugin(ctx)

	// Load dependent external resources
	if err = u.LoadExternal(ctx); err != nil {
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Run tasks
	first, finish, getLastErr := u.RunTask(ctx)

	first.Wait()
	if err = getLastErr(); err != nil {
		cancel()
	}

	// Wait for all tasks to exit
	finish.Wait()
	if err != nil {
		err = getLastErr()
	}

	return
}
