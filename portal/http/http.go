package http

import (
	"context"
	"net/http"
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

// WithHTTPHandler
func WithHTTPHandler(handler http.Handler) ServerOption {
	return newFnServerOption(func(s *Server) { s.srv.Handler = handler })
}

type Server struct {
	srv *http.Server
}

// NewServer
func NewServer(opts ...ServerOption) (s *Server) {
	s = &Server{
		srv: new(http.Server),
	}
	// Apply option
	for _, opt := range opts {
		opt.apply(s)
	}

	return s
}

// Run
func (s *Server) Run(ctx context.Context, addr string) (err error) {
	s.srv.Addr = addr
	// Option: Addr
	if s.srv.Addr == "" {
		s.srv.Addr = ":80"
	}
	go func() {
		<-ctx.Done()
		s.Stop(context.TODO())
	}()
	return s.srv.ListenAndServe()
}

// Stop
func (s *Server) Stop(ctx context.Context) (err error) {
	defer func() {
		s.srv.Close()
	}()
	return s.srv.Shutdown(ctx)
}
