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
		client, used := varp.(*Client)
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
		client, used := varp.(*Client)
		if used {
			client.newOpts = append(client.newOpts,
				grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
			)
		}
		return
	}
}

type Client struct {
	Resource
	newOpts []grpc.DialOption
	cc      *grpc.ClientConn
}

// Init implements Service.
func (client *Client) init() {

}

// Bind implements Service.
func (client *Client) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	for _, setOpt := range opts {
		if !setOpt(client) {
			unused = append(unused, setOpt)
		}
	}
	return
}

// dial implements Service.
func (client *Client) dial(scheme, name string, builder resolver.Builder) (cc *ClientConn, err error) {
	if scheme == "passthrough" {
		client.cc, err = grpc.NewClient(scheme+":///"+name, client.newOpts...)
	} else {
		newOpts := append(client.newOpts, grpc.WithResolvers(builder))
		client.cc, err = grpc.NewClient(scheme+":///"+name, newOpts...)
	}
	return client.cc, err
}

// Dialed implements service.
func (client *Client) Dialed(cc *grpc.ClientConn) {

}
