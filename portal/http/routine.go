package http

import (
	"context"
	"errors"
)

var (
	ErrNoServerRunning = errors.New("no http server is running")
)

type RoutineOption func(r *Routine)

// WithServer
func WithServer(s *Server) RoutineOption {
	return func(r *Routine) { r.s = s }
}

// WithReadyHandler
func WithReadyHandler(fn func(ctx context.Context) error) RoutineOption {
	return func(r *Routine) { r.readyFn = fn }
}

type Routine struct {
	s       *Server
	addr    string
	readyFn func(ctx context.Context) error
}

// NewRoutine
func NewRoutine(addr string, opts ...RoutineOption) (r *Routine) {
	r = &Routine{addr: addr}
	// Set options
	for _, setOpt := range opts {
		setOpt(r)
	}
	// Option: addr
	if r.addr == "" {
		r.addr = ":80"
	}
	// Option: readyFn
	if r.readyFn == nil {
		r.readyFn = func(ctx context.Context) error {
			if r.s != nil {
				return nil
			}
			return ErrNoServerRunning
		}
	}
	return
}

// Name implements core.Routine.
func (r *Routine) Name() string {
	return "NexITF http server"
}

// Run implements core.Routine.
func (r *Routine) Run(ctx context.Context) (err error) {
	if r.s == nil {
		return ErrNoServerRunning
	}
	return r.s.Run(ctx, r.addr)
}

// Stop implements core.Routine.
func (r *Routine) Stop(ctx context.Context) (err error) {
	if r.s == nil {
		return ErrNoServerRunning
	}
	return r.s.Stop(ctx)
}

// Ready implements core.ReadyChecker.
func (h *Routine) Ready(ctx context.Context) (err error) {
	return h.readyFn(ctx)
}
