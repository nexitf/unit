package service

import (
	"net"

	"github.com/nexitf/unit/plugin"
)

type Balancer interface {
	Pick() (addr string, found bool)
	Update(endpoints []Endpoint) (err error)
}

type ConnBalancer interface {
	Pick() (conn net.Conn, found bool)
	Update(endpoints []Endpoint) (err error)
}

// Inner interface: balancer
type balancer interface {
	pick() (addr string, found bool)
	update(endpoints []Endpoint) (err error)
}

// Inner interface: connBalancer
type connBalancer interface {
	pick() (conn net.Conn, found bool)
	update(endpoints []Endpoint) (err error)
}

// WithRoundRobinBalancer binds a RoundRobinBalancer to service.
//
// Support:
//   - HTTPClient
func WithRoundRobinBalancer(endpoints ...Endpoint) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		switch client := varp.(type) {
		// HTTPClient
		case *HTTPClient:
			client.balancer = NewRoundRobinBalancer()
			if len(endpoints) > 0 {
				client.balancer.update(endpoints)
			}
			return true
		}
		return
	}
}

// WithWeightRoundRobinBalancer binds a WeightRoundRobinBalancer to service.
//
// Support:
//   - HTTPClient
func WithWeightRoundRobinBalancer(endpoints ...Endpoint) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		switch client := varp.(type) {
		// HTTPClient
		case *HTTPClient:
			client.balancer = NewWeightRoundRobinBalancer()
			if len(endpoints) > 0 {
				client.balancer.update(endpoints)
			}
			return true
		}
		return
	}
}
