package portal

import (
	"context"
	"net/http"
)

type HTTPServerOption interface {
	apply(*HTTPServer)
}

type fnHTTPServerOption struct {
	f func(*HTTPServer)
}

// newFnHTTPServerOption
func newFnHTTPServerOption(f func(*HTTPServer)) *fnHTTPServerOption {
	return &fnHTTPServerOption{f: f}
}

// apply
func (opt *fnHTTPServerOption) apply(s *HTTPServer) {
	opt.f(s)
}

// WithHandler
func WithHandler(handler http.Handler) HTTPServerOption {
	return newFnHTTPServerOption(func(s *HTTPServer) { s.srv.Handler = handler })
}

type HTTPServer struct {
	srv *http.Server
}

// NewHTTPServer
func NewHTTPServer(opts ...HTTPServerOption) (s *HTTPServer) {
	s = &HTTPServer{
		srv: new(http.Server),
	}
	// Apply option
	for _, opt := range opts {
		opt.apply(s)
	}

	return s
}

// Run
func (s *HTTPServer) Run(ctx context.Context, addr string) (err error) {
	s.srv.Addr = addr
	// Option: Addr
	if s.srv.Addr == "" {
		s.srv.Addr = ":80"
	}
	return s.srv.ListenAndServe()
}

// Stop
func (s *HTTPServer) Stop(ctx context.Context) (err error) {
	defer func() {
		s.srv.Close()
	}()
	return s.srv.Shutdown(ctx)
}
