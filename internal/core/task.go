package core

import (
	"context"
	"sync"
	"sync/atomic"
)

type Task interface {
	Name() string
	Run(ctx context.Context) (err error)
	Ready(ctx context.Context) (err error)
	Stop(ctx context.Context) (err error)
}

// StartTasks
func (u *unit) StartTasks(ctx context.Context) {
	var (
		taskStarted sync.WaitGroup
	)

	u.first.Add(1)
	// Launch tasks
	for _, task := range u.tasks {
		taskStarted.Add(1)
		u.finish.Add(1)
		go func(task Task) {
			taskStarted.Done()
			u.err = task.Run(ctx)
			if u.err == nil {
				task.Stop(ctx)
			}
			u.first.Done()
			u.finish.Done()
		}(task)
	}

	// Wait for all tasks be started
	taskStarted.Wait()

	go func() {
		// Will be blocked here
		u.first.Wait()
		if u.err != nil {
			// Exit all tasks and stop sending readiness notifications
			// when the first runtime error is encountered
			atomic.StoreInt32(&u.noReady, 1)
			u.cancelTask()
		}

		// Wait for all tasks to exit
		u.finish.Wait()
		u.taskExited <- u.err
	}()
}

// ReadyTasks
func (u *unit) ReadyTasks(ctx context.Context) (err error) {
	for _, task := range u.tasks {
		if err = task.Ready(ctx); err != nil {
			return
		}
	}
	return
}
