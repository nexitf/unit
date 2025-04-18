package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
)

type ServerOption interface {
	apply(*Server)
}

type fnServerOption struct {
	f func(*Server)
}

// newFnServerOption
func newFnServerOption(f func(*Server)) *fnServerOption {
	return &fnServerOption{f: f}
}

// apply
func (opt *fnServerOption) apply(s *Server) {
	opt.f(s)
}

// WithServiceRegisterFunc
func WithServiceRegisterFunc(fn func(srv *grpc.Server)) ServerOption {
	return newFnServerOption(func(s *Server) { s.register = fn })
}

type Server struct {
	ln       net.Listener
	srv      *grpc.Server
	register func(srv *grpc.Server)
}

// NewServer
func NewServer(opts ...ServerOption) (s *Server) {
	s = &Server{
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
func (s *Server) Run(ctx context.Context, addr string) (err error) {
	s.ln, err = net.Listen("tcp", addr)
	if err != nil {
		return
	}
	return s.srv.Serve(s.ln)
}

// Stop
func (s *Server) Stop(ctx context.Context) (err error) {
	s.srv.Stop()
	if s.ln != nil {
		return s.ln.Close()
	}
	return
}
