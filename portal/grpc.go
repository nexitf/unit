package portal

import (
	"context"
	"net"

	"google.golang.org/grpc"
)

type GRPCServerOption interface {
	apply(*GRPCServer)
}

type fnGRPCServerOption struct {
	f func(*GRPCServer)
}

// newFnGRPCServerOption
func newFnGRPCServerOption(f func(*GRPCServer)) *fnGRPCServerOption {
	return &fnGRPCServerOption{f: f}
}

// apply
func (opt *fnGRPCServerOption) apply(s *GRPCServer) {
	opt.f(s)
}

// WithServiceRegister
func WithServiceRegister(fn func(srv *grpc.Server)) GRPCServerOption {
	return newFnGRPCServerOption(func(s *GRPCServer) { s.register = fn })
}

type GRPCServer struct {
	srv      *grpc.Server
	ln       net.Listener
	register func(srv *grpc.Server)
}

// NewGRPCServer
func NewGRPCServer(opts ...GRPCServerOption) (s *GRPCServer) {
	s = &GRPCServer{
		srv: grpc.NewServer(),
	}
	// Apply options
	for _, opt := range opts {
		opt.apply(s)
	}

	// Option: Register
	if s.register != nil {
		s.register(s.srv)
	}

	return s
}

// Run
func (s *GRPCServer) Run(ctx context.Context, addr string) (err error) {
	s.ln, err = net.Listen("tcp", addr)
	if err != nil {
		return
	}
	return s.srv.Serve(s.ln)
}

// Stop
func (s *GRPCServer) Stop(ctx context.Context) (err error) {
	s.srv.Stop()
	if s.ln != nil {
		return s.ln.Close()
	}
	return
}
