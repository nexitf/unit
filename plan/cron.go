package plan

import (
	"context"
	"time"

	"github.com/nexitf/logger"
	"github.com/nexitf/logger/field"
	cron "github.com/robfig/cron/v3"
	"go.uber.org/multierr"
)

type CronOption = cron.Option

// WithLocation overrides the timezone of the cron instance.
func WithLocation(loc *time.Location) CronOption {
	return cron.WithLocation(loc)
}

// WithSeconds overrides the parser used for interpreting job schedules to
// include a seconds field as the first one.
func WithSeconds() CronOption {
	return cron.WithSeconds()
}

type CronRoutine struct {
	cron *cron.Cron
	err  error
}

// NewCronRoutine
func NewCronRoutine(opts ...CronOption) (r *CronRoutine) {
	// Set logger option
	opts = append(opts,
		cron.WithLogger(&cronLogger{}),
	)
	r = &CronRoutine{
		cron: cron.New(opts...),
	}
	return
}

// Cmd
func (r *CronRoutine) Cmd(spec string, fn func(time.Time)) {
	_, err := r.cron.AddFunc(spec, func() { fn(time.Now()) })
	r.err = multierr.Append(r.err, err)
}

// Name implements core.Routine.
func (r *CronRoutine) Name() string {
	return "NexITF cron plan"
}

// Run implements core.Routine.
func (r *CronRoutine) Run(ctx context.Context) (err error) {
	if r.err != nil {
		return r.err
	}
	go func() {
		<-ctx.Done()
		r.cron.Stop()
	}()
	r.cron.Run()
	return ctx.Err()
}

// Stop implements core.Routine.
func (r *CronRoutine) Stop(ctx context.Context) (err error) {
	return
}

type cronLogger struct{}

// Info implements cron.Logger.
func (l *cronLogger) Info(msg string, keysAndValues ...any) {
	logger.Info(msg, l.formatFields(keysAndValues...)...)
}

// Error implements cron.Logger.
func (l *cronLogger) Error(err error, msg string, keysAndValues ...any) {
	logger.Error(msg, append([]logger.FieldOption{field.Error(err)}, l.formatFields(keysAndValues...)...)...)
}

// formatFields
func (l *cronLogger) formatFields(keysAndValues ...any) (fields []logger.FieldOption) {
	for i := 0; i < len(keysAndValues); i += 2 {
		if key, ok := keysAndValues[i].(string); ok {
			fields = append(fields, field.Value(key, keysAndValues[i+1]))
		}
	}
	return
}
