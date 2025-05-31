package service

import (
	"github.com/nexitf/unit/plugin"
)

type Balancer interface {
	Pick() (addr string, found bool)
	Update(endpoints []Endpoint) (err error)
}

type BalancerSetter interface {
	SetBalancer(b Balancer)
}

// Inner interface: balancer
type balancer interface {
	pick() (addr string, found bool)
	update(endpoints []Endpoint) (err error)
}

// WithRoundRobinBalancer binds a RoundRobinBalancer to service.
//
// Support:
//   - HTTPClient
//   - BalancerSetter
func WithRoundRobinBalancer() plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		switch client := varp.(type) {
		// HTTPClient
		case *HTTPClient:
			client.balancer = NewRoundRobinBalancer()
			return true
		// BalancerSetter
		case BalancerSetter:
			client.SetBalancer(NewRoundRobinBalancer())
			return true
		}
		return
	}
}

// WithWeightRoundRobinBalancer binds a WeightRoundRobinBalancer to service.
//
// Support:
//   - HTTPClient
//   - BalancerSetter
func WithWeightRoundRobinBalancer() plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		switch client := varp.(type) {
		// HTTPClient
		case *HTTPClient:
			client.balancer = NewWeightRoundRobinBalancer()
			return true
		// BalancerSetter
		case BalancerSetter:
			client.SetBalancer(NewWeightRoundRobinBalancer())
			return true
		}
		return
	}
}
