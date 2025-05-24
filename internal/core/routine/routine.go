package routine

import (
	"context"
	"sync"
	"time"

	"github.com/nexitf/logkit"
	"github.com/nexitf/unit/analytics/stats"
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

	// Stop method will be called once when the routine exits.
	Stop(ctx context.Context) (err error)
}

type RoutineChecker interface {
	Ready(ctx context.Context) (err error)
}

type Launcher struct {
	Routine  Routine
	Err      error
	RunTime  time.Time
	StopTime time.Time
}

// Name implements Routine.
func (l *Launcher) Name() string {
	return l.Routine.Name()
}

// Run implements Routine.
func (l *Launcher) Run(ctx context.Context) (err error) {
	defer func() {
		if l.Err == nil {
			l.Err = err
		}
		l.StopTime = time.Now()
	}()
	l.RunTime = time.Now()
	return l.Routine.Run(ctx)
}

// Stop implements Routine.
func (l *Launcher) Stop(ctx context.Context) (err error) {
	l.StopTime = time.Now()
	return l.Routine.Stop(ctx)
}

// Ready implements RoutineChecker.
func (l *Launcher) Ready(ctx context.Context) (err error) {
	rc, ok := l.Routine.(RoutineChecker)
	if !ok {
		return nil
	}
	defer func() {
		if l.Err == nil {
			l.Err = err
		}
	}()
	return rc.Ready(ctx)
}

func Start(ctx context.Context, routine Routine, started *sync.WaitGroup, report chan error) {
	go func() {
		var err error

		started.Done()
		defer func() { report <- err }()

		l, ok := routine.(*Launcher)
		if !ok {
			report <- ErrRoutineException
			return
		}

		spend := stats.CallSpend(func() {
			err = l.Run(ctx)
		})
		if err != nil {
			logkit.ErrorWrap(err, "routine run failed",
				logkit.Field("routine", l.Name()),
				logkit.Field("spend", spend.Seconds()),
			)
		} else {
			logkit.Info("routine run successfully",
				logkit.Field("routine", l.Name()),
				logkit.Field("spend", spend.Seconds()),
			)
			l.Stop(ctx)
		}
	}()
}
