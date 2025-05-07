package portal

import (
	"context"
	"net"

	"github.com/nexitf/unit/internal/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/stats"
)

var (
	ErrNoGRPCServerRunning = errors.New("no grpc server is running")
)

type GRPCRoutineOption func(r *GRPCRoutine)

// WithWriteBufferSize determines how much data can be batched before doing a write
// on the wire. The default value for this buffer is 32KB. Zero or negative
// values will disable the write buffer such that each write will be on underlying
// connection. Note: A Send call may not directly translate to a write.
func WithWriteBufferSize(size int) GRPCRoutineOption {
	return func(r *GRPCRoutine) { r.grpcOpts = append(r.grpcOpts, grpc.WriteBufferSize(size)) }
}

// WithReadBufferSize lets you set the size of read buffer, this determines how much
// data can be read at most for one read syscall. The default value for this
// buffer is 32KB. Zero or negative values will disable read buffer for a
// connection so data framer can access the underlying conn directly.
func WithReadBufferSize(size int) GRPCRoutineOption {
	return func(r *GRPCRoutine) { r.grpcOpts = append(r.grpcOpts, grpc.ReadBufferSize(size)) }
}

// WithInitialWindowSize sets window size for stream.
// The lower bound for window size is 64K and any value smaller than that will be ignored.
func WithInitialWindowSize(size int32) GRPCRoutineOption {
	return func(r *GRPCRoutine) { r.grpcOpts = append(r.grpcOpts, grpc.InitialWindowSize(size)) }
}

// WithInitialConnWindowSize sets window size for a connection.
// The lower bound for window size is 64K and any value smaller than that will be ignored.
func WithInitialConnWindowSize(size int32) GRPCRoutineOption {
	return func(r *GRPCRoutine) { r.grpcOpts = append(r.grpcOpts, grpc.InitialConnWindowSize(size)) }
}

// WithKeepaliveParams sets keepalive and max-age parameters for the server.
func WithKeepaliveParams(params keepalive.ServerParameters) GRPCRoutineOption {
	return func(r *GRPCRoutine) { r.grpcOpts = append(r.grpcOpts, grpc.KeepaliveParams(params)) }
}

// WithKeepaliveEnforcementPolicy sets keepalive enforcement policy for the server.
func WithKeepaliveEnforcementPolicy(kep keepalive.EnforcementPolicy) GRPCRoutineOption {
	return func(r *GRPCRoutine) { r.grpcOpts = append(r.grpcOpts, grpc.KeepaliveEnforcementPolicy(kep)) }
}

// WithMaxRecvMsgSize sets the max message size in bytes the server can receive.
// If this is not set, gRPC uses the default 4MB.
func WithMaxRecvMsgSize(m int) GRPCRoutineOption {
	return func(r *GRPCRoutine) { r.grpcOpts = append(r.grpcOpts, grpc.MaxRecvMsgSize(m)) }
}

// WithMaxSendMsgSize sets the max message size in bytes the server can send.
// If this is not set, gRPC uses the default `math.MaxInt32`.
func WithMaxSendMsgSize(m int) GRPCRoutineOption {
	return func(r *GRPCRoutine) { r.grpcOpts = append(r.grpcOpts, grpc.MaxSendMsgSize(m)) }
}

// WithMaxConcurrentStreams applies a limit on the number
// of concurrent streams to each ServerTransport.
func WithMaxConcurrentStreams(n uint32) GRPCRoutineOption {
	return func(r *GRPCRoutine) { r.grpcOpts = append(r.grpcOpts, grpc.MaxConcurrentStreams(n)) }
}

// WithCreds sets credentials for server connections.
func WithCreds(creds credentials.TransportCredentials) GRPCRoutineOption {
	return func(r *GRPCRoutine) { r.grpcOpts = append(r.grpcOpts, grpc.Creds(creds)) }
}

// WithUnaryInterceptors specifies the chained interceptor
// for unary RPCs. The first interceptor will be the outer most,
// while the last interceptor will be the inner most wrapper around the real call.
// All unary interceptors added by this method will be chained.
func WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) GRPCRoutineOption {
	return func(r *GRPCRoutine) { r.grpcOpts = append(r.grpcOpts, grpc.ChainUnaryInterceptor(interceptors...)) }
}

// WithStreamInterceptors specifies the chained interceptor
// for streaming RPCs. The first interceptor will be the outer most,
// while the last interceptor will be the inner most wrapper around the real call.
// All stream interceptors added by this method will be chained.
func WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) GRPCRoutineOption {
	return func(r *GRPCRoutine) { r.grpcOpts = append(r.grpcOpts, grpc.ChainStreamInterceptor(interceptors...)) }
}

// WithStatsHandler sets the stats handler for the server.
func WithStatsHandler(handler stats.Handler) GRPCRoutineOption {
	return func(r *GRPCRoutine) { r.grpcOpts = append(r.grpcOpts, grpc.StatsHandler(handler)) }
}

// WithUnknownServiceHandler allows for adding a custom
// unknown service handler. The provided method is a bidi-streaming RPC service
// handler that will be invoked instead of returning the "unimplemented" gRPC
// error whenever a request is received for an unregistered service or method.
// The handling function and stream interceptor (if set) have full access to
// the ServerStream, including its Context.
func WithUnknownServiceHandler(handler grpc.StreamHandler) GRPCRoutineOption {
	return func(r *GRPCRoutine) { r.grpcOpts = append(r.grpcOpts, grpc.UnknownServiceHandler(handler)) }
}

// WithMaxHeaderListSize sets the max (uncompressed) size
// of header list that the server is prepared to accept.
func WithMaxHeaderListSize(size uint32) GRPCRoutineOption {
	return func(r *GRPCRoutine) { r.grpcOpts = append(r.grpcOpts, grpc.MaxHeaderListSize(size)) }
}

// WithGRPCServiceRegisterHandler
func WithGRPCServiceRegisterHandler(fn func(srv *grpc.Server)) GRPCRoutineOption {
	return func(r *GRPCRoutine) { r.regFn = fn }
}

// WithGRPCReadyHandler
func WithGRPCReadyHandler(fn func(ctx context.Context) error) GRPCRoutineOption {
	return func(r *GRPCRoutine) { r.readyFn = fn }
}

type GRPCRoutine struct {
	srv      *grpc.Server
	grpcOpts []grpc.ServerOption
	ln       net.Listener
	addr     string
	regFn    func(srv *grpc.Server)
	readyFn  func(ctx context.Context) error
}

// NewGRPCRoutine
func NewGRPCRoutine(addr string, opts ...GRPCRoutineOption) (r *GRPCRoutine) {
	r = &GRPCRoutine{
		addr: addr,
	}
	// Set options
	for _, setOpt := range opts {
		setOpt(r)
	}
	// New grpc server
	r.srv = grpc.NewServer(r.grpcOpts...)
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
