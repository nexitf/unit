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
)

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

// WithRoundRobinBalancer sets Round-Robin Balancer.
func WithRoundRobinBalancer() plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy": "`+roundrobin.Name+`"}`),
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
				grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy": "`+weightedroundrobin.Name+`"}`),
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
