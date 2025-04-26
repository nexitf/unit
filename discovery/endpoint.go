package discovery

import (
	"sort"
)

type Endpoint struct {
	ID     int    `json:"id,omitempty"`
	Addr   string `json:"addr,omitempty"`
	Weight int    `json:"weight,omitempty"`
	Meta   string `json:"meta,omitempty"`
}

func SortByID(endpoints []Endpoint) []Endpoint {
	sort.Slice(endpoints, func(i, j int) bool { return endpoints[i].ID < endpoints[j].ID })
	return endpoints
}

func SortByAddr(endpoints []Endpoint) []Endpoint {
	sort.Slice(endpoints, func(i, j int) bool { return endpoints[i].Addr < endpoints[j].Addr })
	return endpoints
}
