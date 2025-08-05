package routine

import (
	"context"
	"sync"
	"time"

	"github.com/nexitf/logger"
	"github.com/nexitf/logger/field"
	"github.com/nexitf/unit/analytics/stats"
	"github.com/nexitf/unit/internal/core/utils"
)

type Launcher struct {
	routine  Routine
	err      error
	runTime  time.Time
	stopTime time.Time
	mutex    sync.RWMutex
}

func (l *Launcher) Name() string {
	return l.routine.Name()
}

// Start initiates the wrapped routine in a new goroutine and monitors its execution.
// It records the execution time of the routine, logs the result based on its success or failure,
// and stops the routine if it executes successfully.
//
// Parameters:
//   - ctx: The context used to control the lifecycle of the routine, which can propagate cancellation signals.
//   - started: A sync.WaitGroup that notifies the caller when the routine has started.
//   - report: A channel for errors, used to pass any errors encountered during the routine's execution.
func (l *Launcher) Start(ctx context.Context, started *sync.WaitGroup, report chan error) {
	go func() {
		var err error

		started.Done()
		defer func() { report <- err }()

		spend := stats.CallSpend(func() {
			err = l.Run(ctx)
		})
		if err != nil {
			logger.Error("routine run failed",
				field.Error(err),
				field.Value("routine", l.Name()),
				field.Value("spend", spend.Seconds()),
			)
		} else {
			logger.Info("routine run successfully",
				field.Value("routine", l.Name()),
				field.Value("spend", spend.Seconds()),
			)
			l.Stop(ctx)
		}
	}()
}

// Run executes the wrapped routine and records its start and stop times, as well as any errors encountered.
// It ensures that the routine's error and stop time are recorded even if the routine panics.
//
// Parameters:
//   - ctx: The context used to control the lifecycle of the routine, which can propagate cancellation signals.
//
// Returns:
//   - err: The error returned by the wrapped routine, or nil if the routine executed successfully.
func (l *Launcher) Run(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			l.reason(err)
		}
		l.mutex.Lock()
		l.stopTime = time.Now()
		l.mutex.Unlock()
	}()
	l.mutex.Lock()
	l.runTime = time.Now()
	l.mutex.Unlock()
	return l.routine.Run(ctx)
}

// Stop terminates the wrapped routine and records the stop time.
// It calls the Stop method of the underlying Routine interface, passing the provided context.
// The stop time is recorded immediately before attempting to stop the routine.
//
// Parameters:
//   - ctx: The context used to control the termination process, which can propagate cancellation signals.
//
// Returns:
//   - err: The error returned by the underlying routine's Stop method, or nil if the routine stopped successfully.
func (l *Launcher) Stop(ctx context.Context) (err error) {
	r, ok := l.routine.(RoutineStopped)
	if !ok {
		return nil
	}
	return r.Stop(ctx)
}

// Ready checks if the wrapped routine is ready.
// If the wrapped routine implements the RoutineChecker interface, it calls the Ready method of that interface to perform the readiness check.
// If the routine does not implement the RoutineChecker interface, it returns nil directly, indicating that it is considered ready by default.
// Regardless of the check result, any error will be recorded in the Err field of the Launcher.
//
// Parameters:
//   - ctx: The context used to control the checking process, which can propagate cancellation signals.
//
// Returns:
//   - err: The error encountered during the readiness check. Returns nil if the check succeeds or the routine does not implement the RoutineChecker interface.
func (l *Launcher) Ready(ctx context.Context) (err error) {
	rc, ok := l.routine.(RoutineChecker)
	if !ok {
		return nil
	}
	defer func() {
		if err != nil {
			l.reason(err)
		}
	}()
	return rc.Ready(ctx)
}

// RunTime returns the runtime duration of the routine.
func (l *Launcher) RunTime() time.Time {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.runTime
}

// StopTime returns the time when the routine stopped.
func (l *Launcher) StopTime() time.Time {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.stopTime
}

// Err returns the error occurred during the routine's execution.
func (l *Launcher) Err() (err error) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.err
}

// reason sets the error encountered during the execution of routines,
// If an error has already been set, it does nothing.
func (l *Launcher) reason(err error) {
	if l.err != nil {
		return
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.err == nil {
		l.err = err
	}
}

// Group implements the Routine interface and is used to set routines in batches.
// Group itself will not be executed as a routine.
type Group struct {
	routines []Routine
}

// NewGroup creates a new instance of Group.
// It initializes an empty slice of routines and returns a pointer to the newly created Group.
// This function serves as a convenient way to instantiate a Group.
func NewGroup() (g *Group) {
	return &Group{}
}

// Name returns the name of the Group.
// This method implements the corresponding method in the routine.Routine interface,
// providing a standardized way to identify the Group instance.
// It always returns the fixed string "NEXITF routine group".
//
// Returns:
// - string: The name of the Group.
func (g *Group) Name() string {
	return "NEXITF routine group"
}

// Add appends a new routine to the Group.
// It takes a routine that implements the routine.Routine interface as a parameter.
// The provided routine is then added to the internal slice of routines within the Group.
// This method modifies the state of the Group by expanding its list of managed routines.
func (g *Group) Add(routine Routine) {
	g.routines = append(g.routines, routine)
}

// Run will never be executed.
func (g *Group) Run(ctx context.Context) (err error) {
	return
}

type Scheduler struct {
	routines  []*Launcher   // A slice of Launcher instances, each wrapping a Routine.
	stopped   chan struct{} // A channel that is closed when all routines have exited.
	err       error         // The reason for the scheduler's termination, if any.
	mutex     sync.RWMutex
	cancelCtx context.CancelFunc
}

// NewScheduler creates a new instance of Scheduler.
// It initializes the stopped channel to nil, which will be properly initialized when the scheduler starts.
func NewScheduler() (s *Scheduler) {
	return &Scheduler{stopped: nil}
}

// Routines returns a slice of pointers to the Launcher instances managed by the scheduler.
// This method provides direct access to the internal list of routines, which can be used
// to inspect the state of each routine. Note that modifying the returned slice or the
// Launcher instances it contains may affect the scheduler's internal state.
//
// Returns:
//   - A slice of pointers to Launcher instances.
func (s *Scheduler) Routines() []*Launcher {
	return s.routines
}

// Num returns the number of routines currently managed by the scheduler.
// It can be used to get an overview of how many routines are part of the scheduler.
// Returns:
//   - An integer representing the count of routines in the scheduler.
func (s *Scheduler) Num() int {
	return len(s.routines)
}

// Add adds a routine or group to the scheduler.
// It wraps the provided Routine in a Launcher and appends it to the list of routines.
func (s *Scheduler) Add(routine Routine) {
	if g, ok := routine.(*Group); !ok {
		s.routines = append(s.routines, &Launcher{routine: routine})
	} else {
		for _, r := range g.routines {
			s.routines = append(s.routines, &Launcher{routine: r})
		}
	}
}

// Start launches all the routines managed by the scheduler.
// It creates a cancellable context, starts each routine, and waits for them to be started.
// If any routine returns an error, it cancels the context and records the first error.
// It also listens for all routines to exit and closes the stopped channel when done.
func (s *Scheduler) Start(ctx context.Context) {
	var (
		started sync.WaitGroup
		num     = len(s.routines)
		report  = make(chan error, num)
	)

	// Initialize the stopped channel to signal when all routines have exited.
	s.stopped = make(chan struct{})
	ctx, s.cancelCtx = context.WithCancel(ctx)

	// Launch all routines
	for _, r := range s.routines {
		started.Add(1)
		r.Start(ctx, &started, report)
	}

	// Wait for all routines be started
	started.Wait()

	go func() {
		// Will be blocked here
		for err := range report {
			num--
			if err != nil {
				s.Interrupt(err)
			}

			// When all routines exited
			if num <= 0 {
				close(report)
				break
			}
		}

		// Exit
		close(s.stopped)
	}()
}

func (s *Scheduler) Ready(ctx context.Context) {
	if s.Err() != nil {
		return
	}
	// Check ready
	for _, r := range s.routines {
		var err error
		spend := utils.CallSpend(func() {
			err = r.Ready(ctx)
		})
		if err != nil {
			logger.Error("routine readiness check failed",
				field.Error(err),
				field.Value("routine", r.Name()),
				field.Value("spend", spend.Seconds()),
			)
			s.Interrupt(err)
			return
		}
		logger.Info("routine readiness check successfully",
			field.Value("routine", r.Name()),
			field.Value("spend", spend.Seconds()),
		)
	}
}

// Interrupt is used to interrupt the execution of all routines managed by the scheduler.
// It records the provided error as the reason for the interruption and cancels the context
// associated with all routines, which will signal them to stop their execution.
//
// Parameters:
//   - err: The error that serves as the reason for interrupting the routines.
//     This error will be stored and can be retrieved later using the Err method.
func (s *Scheduler) Interrupt(err error) {
	s.reason(err)
	s.cancelCtx()
}

// Done returns a read-only channel that is closed when all routines have exited.
// It can be used to wait for the completion of all routines.
// Returns:
//   - A read-only channel of struct{}.
func (s *Scheduler) Done() <-chan struct{} {
	return s.stopped
}

// Err returns the first error encountered during the execution of routines.
// If no errors occurred, it returns nil.
// Returns:
//   - The first error encountered, or nil if none.
func (s *Scheduler) Err() (err error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.err
}

// reason sets the first error encountered during the execution of routines.
// If an error has already been set, it does nothing.
// This method is thread-safe, using a mutex to ensure that only one goroutine can set the error at a time.
//
// Parameters:
//   - err: The error to be set as the reason for the scheduler's termination.
func (s *Scheduler) reason(err error) {
	if s.err != nil {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.err == nil {
		s.err = err
	}
}
