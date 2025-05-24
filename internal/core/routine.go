package core

import (
	"context"
	"sync"

	"github.com/nexitf/logkit"
	"github.com/nexitf/unit/internal/core/routine"
	"github.com/nexitf/unit/internal/core/utils"
)

// StartRoutines
func (u *Unit) StartRoutines(ctx context.Context) {
	var (
		started sync.WaitGroup
		num     = len(u.routines)
		report  = make(chan error, num)
	)

	// Launch all routines
	for _, r := range u.routines {
		started.Add(1)
		routine.Start(ctx, r, &started, report)
	}

	// Wait for all routines be started
	started.Wait()

	go func() {
		// Will be blocked here
		for err := range report {
			num--
			if err != nil {
				u.Panic(err)
				u.cancelCtx()
			}
			// When all routines exited
			if num <= 0 {
				close(report)
				break
			}
		}

		// Exit
		close(u.unitExited)
	}()
}

// ReadyRoutines
func (u *Unit) ReadyRoutines(ctx context.Context) (err error) {
	for _, r := range u.routines {
		rc, ok := r.(routine.RoutineChecker)
		if !ok {
			continue
		}
		spend := utils.CallSpend(func() {
			err = rc.Ready(ctx)
		})
		if err != nil {
			logkit.ErrorWrap(err, "routine readiness check failed",
				logkit.Field("routine", r.Name()),
				logkit.Field("spend", spend.Seconds()),
			)
			u.Panic(err)
			return
		}
		logkit.Info("routine readiness check successfully",
			logkit.Field("routine", r.Name()),
			logkit.Field("spend", spend.Seconds()),
		)
	}
	return
}
