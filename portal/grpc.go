package portal

import (
	"context"
	"net"

	"github.com/nexitf/unit/internal/errors"
	"google.golang.org/grpc"
)

var (
	ErrNoGRPCServerRunning = errors.New("no grpc server is running")
)

type GRPCRoutineOption func(r *GRPCRoutine)

// WithGRPCServiceRegisterHandler
func WithGRPCServiceRegisterHandler(fn func(srv *grpc.Server)) GRPCRoutineOption {
	return func(r *GRPCRoutine) { r.regFn = fn }
}

// WithGRPCReadyHandler
func WithGRPCReadyHandler(fn func(ctx context.Context) error) GRPCRoutineOption {
	return func(r *GRPCRoutine) { r.readyFn = fn }
}

type GRPCRoutine struct {
	srv     *grpc.Server
	ln      net.Listener
	addr    string
	regFn   func(srv *grpc.Server)
	readyFn func(ctx context.Context) error
}

// NewGRPCRoutine
func NewGRPCRoutine(addr string, opts ...GRPCRoutineOption) (r *GRPCRoutine) {
	r = &GRPCRoutine{
		srv:  grpc.NewServer(),
		addr: addr,
	}
	// Set options
	for _, setOpt := range opts {
		setOpt(r)
	}
	// Option: addr
	if r.addr == "" {
		r.addr = ":8999"
	}
	// Option: regFn
	if r.regFn != nil {
		r.regFn(r.srv)
	}
	// Option: readyFn
	if r.readyFn == nil {
		r.readyFn = func(ctx context.Context) error {
			if r.srv != nil {
				return nil
			}
			return ErrNoGRPCServerRunning
		}
	}
	return
}

// Name implements core.Routine.
func (r *GRPCRoutine) Name() string {
	return "NexITF grpc server"
}

// Run implements core.Routine.
func (r *GRPCRoutine) Run(ctx context.Context) (err error) {
	if r.srv == nil {
		return ErrNoGRPCServerRunning
	}
	r.ln, err = net.Listen("tcp", r.addr)
	if err != nil {
		return
	}
	go func() {
		<-ctx.Done()
		r.srv.Stop()
	}()
	return r.srv.Serve(r.ln)
}

// Stop implements core.Routine.
func (r *GRPCRoutine) Stop(ctx context.Context) (err error) {
	if r.srv == nil {
		return ErrNoGRPCServerRunning
	}
	r.srv.Stop()
	if r.ln != nil {
		return r.ln.Close()
	}
	return
}

// Ready implements core.ReadyChecker.
func (h *GRPCRoutine) Ready(ctx context.Context) (err error) {
	return h.readyFn(ctx)
}
