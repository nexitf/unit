package discovery

import (
	"context"
)

var (
	// An empty implementation of Service Discovery
	// that never triggers any updates.
	Offline = &emptyDiscovery{}
)

type emptyDiscovery struct {
}

// Watch
func (*emptyDiscovery) Watch(ctx context.Context, serviceName, tag string, update func(endpoints []Endpoint, closed bool)) (cancel func() error, err error) {
	return func() error { return nil }, nil
}
