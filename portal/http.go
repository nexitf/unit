package portal

import (
	"context"
	"net/http"

	"github.com/nexitf/unit/internal/errors"
)

var (
	ErrNoHTTPServerRunning = errors.New("no http server is running")
)

type HTTPRoutineOption func(r *HTTPRoutine)

// WithHTTPHandler
func WithHTTPHandler(handler http.Handler) HTTPRoutineOption {
	return func(r *HTTPRoutine) { r.srv.Handler = handler }
}

// WithHTTPReadyHandler
func WithHTTPReadyHandler(fn func(ctx context.Context) error) HTTPRoutineOption {
	return func(r *HTTPRoutine) { r.readyFn = fn }
}

type HTTPRoutine struct {
	srv     *http.Server
	addr    string
	readyFn func(ctx context.Context) error
}

// NewHTTPRoutine
func NewHTTPRoutine(addr string, opts ...HTTPRoutineOption) (r *HTTPRoutine) {
	r = &HTTPRoutine{
		srv:  new(http.Server),
		addr: addr,
	}
	// Set options
	for _, setOpt := range opts {
		setOpt(r)
	}
	// Option: addr
	if r.addr == "" {
		r.addr = ":80"
	}
	r.srv.Addr = r.addr
	// Option: readyFn
	if r.readyFn == nil {
		r.readyFn = func(ctx context.Context) error {
			if r.srv != nil {
				return nil
			}
			return ErrNoHTTPServerRunning
		}
	}
	return
}

// Name implements core.Routine.
func (r *HTTPRoutine) Name() string {
	return "NexITF http server"
}

// Run implements core.Routine.
func (r *HTTPRoutine) Run(ctx context.Context) (err error) {
	if r.srv == nil {
		return ErrNoHTTPServerRunning
	}
	go func() {
		<-ctx.Done()
		r.srv.Shutdown(ctx)
	}()
	return r.srv.ListenAndServe()
}

// Stop implements core.Routine.
func (r *HTTPRoutine) Stop(ctx context.Context) (err error) {
	if r.srv == nil {
		return ErrNoHTTPServerRunning
	}
	return r.srv.Close()
}

// Ready implements core.ReadyChecker.
func (h *HTTPRoutine) Ready(ctx context.Context) (err error) {
	return h.readyFn(ctx)
}
