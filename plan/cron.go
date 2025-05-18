package plan

import (
	"context"
	"time"

	"github.com/nexitf/logkit"
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
		cron.WithLogger(&logger{}),
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

type logger struct{}

// Info implements cron.Logger.
func (l *logger) Info(msg string, keysAndValues ...interface{}) {
	logkit.Info(msg, l.formatFields(keysAndValues...)...)
}

// Error implements cron.Logger.
func (l *logger) Error(err error, msg string, keysAndValues ...interface{}) {
	logkit.ErrorWrap(err, msg, l.formatFields(keysAndValues...)...)
}

// formatFields
func (l *logger) formatFields(keysAndValues ...interface{}) (fields []logkit.FieldOption) {
	for i := 0; i < len(keysAndValues); i += 2 {
		if key, ok := keysAndValues[i].(string); ok {
			fields = append(fields, logkit.Field(key, keysAndValues[i+1]))
		}
	}
	return
}
