package grpc

import (
	"github.com/nexitf/unit/plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

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

// WithRoundRobinBalancer sets Round-Robin Balancer.
func WithRoundRobinBalancer() plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*Base)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
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

// Dialed implements service.
func (base *Base) Dialed(cc *grpc.ClientConn) {

}
