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
func WithRoundRobinBalancer() plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		switch client := varp.(type) {
		// HTTPClient
		case *HTTPClient:
			client.balancer = NewRoundRobinBalancer()
			return true
		}
		return
	}
}

// WithWeightRoundRobinBalancer binds a WeightRoundRobinBalancer to service.
//
// Support:
//   - HTTPClient
func WithWeightRoundRobinBalancer() plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		switch client := varp.(type) {
		// HTTPClient
		case *HTTPClient:
			client.balancer = NewWeightRoundRobinBalancer()
			return true
		}
		return
	}
}
