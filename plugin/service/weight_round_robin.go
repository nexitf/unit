package service

import (
	"sync"
)

type node struct {
	current   int
	effective int
	endpoint  Endpoint
}

type WeightRoundRobinBalancer struct {
	rss []*node
	rsm map[string]int
	mu  sync.RWMutex
}

// NewWeightRoundRobinBalancer
func NewWeightRoundRobinBalancer() (b *WeightRoundRobinBalancer) {
	return &WeightRoundRobinBalancer{
		rss: make([]*node, 0),
		rsm: make(map[string]int),
	}
}

// Pick implements Balancer.
func (b *WeightRoundRobinBalancer) Pick() (addr string, found bool) {
	return b.pick()
}

// Update implements Balancer.
func (b *WeightRoundRobinBalancer) Update(endpoints []Endpoint) (err error) {
	return b.update(endpoints)
}

// // UnStable
// func (b *WeightRoundRobinBalancer) UnStable(addr string) {
// 	b.mu.Lock()
// 	defer b.mu.Unlock()
// 	if node, ok := b.node(addr); ok && node.effective > 0 {
// 		node.effective = (int)(node.effective / 2)
// 	}
// }

// Add implements Balancer.
func (b *WeightRoundRobinBalancer) Add(endpoints ...Endpoint) {
	b.mu.Lock()
	defer b.mu.Unlock()
	// Add
	for _, node := range b.rss {
		endpoints = append(endpoints, node.endpoint)
	}
	// Update
	b.update(endpoints)
}

// Remove implements Balancer.
func (b *WeightRoundRobinBalancer) Remove(addr string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	// Remove
	if i, ok := b.rsm[addr]; !ok {
		return
	} else if i < len(b.rss) {
		if i+1 == len(b.rss) {
			b.rss = b.rss[0:i]
		} else {
			b.rss = append(b.rss[0:i], b.rss[i+1:]...)
			for j := i; j < len(b.rss); i++ {
				b.rsm[b.rss[j].endpoint.Addr] = j
			}
		}
		delete(b.rsm, addr)
	}
}

// IsDeprecated
func (b *WeightRoundRobinBalancer) IsDeprecated(addr string) (yes bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.node(addr)
	return !ok
}

// pick implements balancer.
func (b *WeightRoundRobinBalancer) pick() (addr string, found bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	var (
		total int
		best  *node
	)
	for i := 0; i < len(b.rss); i++ {
		node := b.rss[i]
		total += node.effective
		node.current += node.effective
		// -1 when the connection is abnormal,
		// +1 when the communication is successful
		if node.effective < node.endpoint.Weight {
			node.effective++
		}
		// Maximum temporary weight node
		if best == nil || node.current > best.current {
			best = node
		}
	}
	if best != nil {
		best.current -= total
		addr, found = best.endpoint.Addr, true
	}
	return
}

// update implements balancer.
func (b *WeightRoundRobinBalancer) update(endpoints []Endpoint) (err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	// Update
	var (
		rss = make([]*node, 0)
		rsm = make(map[string]int)
	)
	for _, endpoint := range endpoints {
		node := &node{
			endpoint:  endpoint,
			effective: endpoint.Weight,
		}
		if i, ok := b.rsm[endpoint.Addr]; ok {
			node.current = b.rss[i].current
			node.effective = b.rss[i].effective
		}
		rss = append(rss, node)
	}
	for i, node := range rss {
		rsm[node.endpoint.Addr] = i
	}
	// Replace
	b.rss, b.rsm = rss, rsm
	return
}

// node
func (b *WeightRoundRobinBalancer) node(addr string) (node *node, ok bool) {
	if i, ok1 := b.rsm[addr]; !ok1 {
		return nil, false
	} else if i < len(b.rss) {
		node, ok = b.rss[i], true
	}
	return
}
