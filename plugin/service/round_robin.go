package service

import (
	"sync"
	"sync/atomic"
)

type RoundRobinBalancer struct {
	mutex   sync.RWMutex
	counter uint32
	rss     []Endpoint
	rsm     map[string]int
}

// NewRoundRobinBalancer
func NewRoundRobinBalancer() (b *RoundRobinBalancer) {
	return &RoundRobinBalancer{
		rss: make([]Endpoint, 0),
		rsm: make(map[string]int),
	}
}

// Pick implements Balancer.
func (b *RoundRobinBalancer) Pick() (addr string, found bool) {
	return b.pick()
}

// pick implements balancer.
func (b *RoundRobinBalancer) pick() (addr string, found bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if found = len(b.rss) > 0; !found {
		return "", false
	}
	addr = b.rss[int(atomic.AddUint32(&b.counter, 1))%len(b.rss)].Addr
	return
}

// Update implements Balancer.
func (b *RoundRobinBalancer) Update(endpoints []Endpoint) (err error) {
	return b.update(endpoints)
}

// Add adds some new endpoints.
func (b *RoundRobinBalancer) Add(endpoints ...Endpoint) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	// Update
	b.update(append(endpoints, b.rss...))
}

// Remove removes the endpoint.
func (b *RoundRobinBalancer) Remove(addr string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	// Remove
	if i, ok := b.rsm[addr]; !ok {
		return
	} else if i < len(b.rss) {
		if i+1 == len(b.rss) {
			b.rss = b.rss[0:i]
		} else {
			b.rss = append(b.rss[0:i], b.rss[i+1:]...)
			for j := i; j < len(b.rss); i++ {
				b.rsm[b.rss[j].Addr] = j
			}
		}
		delete(b.rsm, addr)
	}
}

// IsDeprecated check if the addr is deprecated.
func (b *RoundRobinBalancer) IsDeprecated(addr string) (yes bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	_, ok := b.rsm[addr]
	return !ok
}

// update implements balancer.
func (b *RoundRobinBalancer) update(endpoints []Endpoint) (err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	var (
		rsm = make(map[string]int)
		rss = make([]Endpoint, len(endpoints))
	)
	for i, endpoint := range endpoints {
		rss[i], rsm[endpoint.Addr] = endpoint, i
	}
	b.rss, b.rsm = rss, rsm
	return
}
