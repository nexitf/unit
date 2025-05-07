package grpc

import (
	"context"
	"net"

	"github.com/nexitf/unit/plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/balancer/weightedroundrobin"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/stats"
)

// WithWriteBufferSize determines how much data can be batched before doing a
// write on the wire. The default value for this buffer is 32KB.
//
// Zero or negative values will disable the write buffer such that each write
// will be on underlying connection. Note: A Send call may not directly
// translate to a write.
func WithWriteBufferSize(size int) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithWriteBufferSize(size),
			)
		}
		return
	}
}

// WithReadBufferSize lets you set the size of read buffer, this determines how
// much data can be read at most for each read syscall.
//
// The default value for this buffer is 32KB. Zero or negative values will
// disable read buffer for a connection so data framer can access the
// underlying conn directly.
func WithReadBufferSize(size int) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithReadBufferSize(size),
			)
		}
		return
	}
}

// WithInitialWindowSize sets the value for initial
// window size on a stream. The lower bound for window size is 64K and any value
// smaller than that will be ignored.
func WithInitialWindowSize(size int32) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithInitialWindowSize(size),
			)
		}
		return
	}
}

// WithInitialConnWindowSize sets the value for
// initial window size on a connection. The lower bound for window size is 64K
// and any value smaller than that will be ignored.
func WithInitialConnWindowSize(size int32) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithInitialConnWindowSize(size),
			)
		}
		return
	}
}

// WithDefaultCallOptions sets the default
// CallOptions for calls over the connection.
func WithDefaultCallOptions(opts ...grpc.CallOption) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithDefaultCallOptions(opts...),
			)
		}
		return
	}
}

// WithConnectParams configures the ClientConn to use the provided
// ConnectParams for creating and maintaining connections to servers.
func WithConnectParams(params grpc.ConnectParams) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithConnectParams(params),
			)
		}
		return
	}
}

// WithInsecure uses unencrypted network communication.
func WithInsecure() plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
		}
		return
	}
}

// WithTransportCredentials configures a connection level security credentials (e.g., TLS/SSL).
// This should not be used together with WithCredentialsBundle.
func WithTransportCredentials(creds credentials.TransportCredentials) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithTransportCredentials(creds),
			)
		}
		return
	}
}

// WithPerRPCCredentials sets credentials and places auth state on each outbound RPC.
func WithPerRPCCredentials(creds credentials.PerRPCCredentials) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithPerRPCCredentials(creds),
			)
		}
		return
	}
}

// WithContextDialer sets a dialer to create connections.
func WithContextDialer(fn func(context.Context, string) (net.Conn, error)) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithContextDialer(fn),
			)
		}
		return
	}
}

// WithStatsHandler specifies the stats handler for
// all the RPCs and underlying network connections in this ClientConn.
func WithStatsHandler(handler stats.Handler) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithStatsHandler(handler),
			)
		}
		return
	}
}

// WithUserAgent specifies a user agent for all the RPCs.
func WithUserAgent(ua string) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithUserAgent(ua),
			)
		}
		return
	}
}

// WithKeepaliveParams specifies keepalive parameters for the client transport.
func WithKeepaliveParams(params keepalive.ClientParameters) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithKeepaliveParams(params),
			)
		}
		return
	}
}

// WithUnaryInterceptors specifies the chained
// interceptor for unary RPCs. The first interceptor will be the outer most,
// while the last interceptor will be the inner most wrapper around the real call.
// All interceptors added by this method will be chained, and the interceptor
// defined by WithUnaryInterceptor will always be prepended to the chain.
func WithUnaryInterceptors(interceptors ...grpc.UnaryClientInterceptor) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithChainUnaryInterceptor(interceptors...),
			)
		}
		return
	}
}

// WithStreamInterceptors specifies the chained
// interceptor for streaming RPCs. The first interceptor will be the outer most,
// while the last interceptor will be the inner most wrapper around the real call.
// All interceptors added by this method will be chained, and the interceptor
// defined by WithStreamInterceptor will always be prepended to the chain.
func WithStreamInterceptors(interceptors ...grpc.StreamClientInterceptor) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithChainStreamInterceptor(interceptors...),
			)
		}
		return
	}
}

// WithAuthority specifies the value to be used as the :authority pseudo-header
// and as the server name in authentication handshake.
func WithAuthority(a string) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithAuthority(a),
			)
		}
		return
	}
}

// WithDisableServiceConfig causes gRPC to ignore any
// service config provided by the resolver and provides a hint to the resolver
// to not fetch service configs.
//
// Note that this dial option only disables service config from resolver. If
// default service config is provided, gRPC will use the default service config.
func WithDisableServiceConfig() plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithDisableServiceConfig(),
			)
		}
		return
	}
}

// WithRoundRobinBalancer sets Round-Robin Balancer.
func WithRoundRobinBalancer() plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"`+roundrobin.Name+`":{}}]}`),
			)
		}
		return
	}
}

// WithWeightedRoundRobinBalancer sets Weighted-Round-Robin Balancer.
func WithWeightedRoundRobinBalancer() plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"`+weightedroundrobin.Name+`":{}}]}`),
			)
		}
		return
	}
}

// WithDisableRetry disables retries.
func WithDisableRetry(ua string) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithDisableRetry(),
			)
		}
		return
	}
}

// WithMaxHeaderListSize specifies the maximum
// (uncompressed) size of header list that the client is prepared to accept.
func WithMaxHeaderListSize(size uint32) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithMaxHeaderListSize(size),
			)
		}
		return
	}
}

type Base struct {
	Resource
	newOpts []grpc.DialOption
	cc      *grpc.ClientConn
}

// Init implements client.
func (base *Base) init() {

}

// Bind implements client.
func (base *Base) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	for _, setOpt := range opts {
		if !setOpt(base) {
			unused = append(unused, setOpt)
		}
	}
	return
}

// dial implements client.
func (base *Base) dial(scheme, name string, builder resolver.Builder) (cc *ClientConn, err error) {
	if scheme == "passthrough" {
		base.cc, err = grpc.NewClient(scheme+":///"+name, base.newOpts...)
	} else {
		newOpts := append(base.newOpts, grpc.WithResolvers(builder))
		base.cc, err = grpc.NewClient(scheme+":///"+name, newOpts...)
	}
	return base.cc, err
}

// OnConnect implements client.
func (base *Base) OnConnect(cc *grpc.ClientConn) {

}

// close implements client.
func (base *Base) close() (err error) {
	if base.cc != nil {
		err = base.cc.Close()
	}
	return
}
