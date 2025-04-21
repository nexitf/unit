package service

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/nexitf/lamp"
	"github.com/nexitf/unit/plugin"
)

var (
	ErrAddressNotFound = errors.New("http address not found")
)

// WithHTTPHost
func WithHTTPHost(host string) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*HTTPClient)
		if used {
			client.host = host
		}
		return
	}
}

type HTTPClient struct {
	serviceBase
	mutex sync.RWMutex
	index uint32
	addrs []lamp.Address
	host  string
}

// Get
func (client *HTTPClient) Get(path string) (body string, err error) {
	addr, found := client.pick()
	if !found {
		return "", ErrAddressNotFound
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	resp, err := http.Get(fmt.Sprintf("http://%s%s", addr.Addr, path))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// pick
func (client *HTTPClient) pick() (addr lamp.Address, found bool) {
	client.mutex.RLock()
	defer client.mutex.RUnlock()
	// No address
	found = len(client.addrs) > 0
	if !found {
		return
	}
	addr = client.addrs[atomic.AddUint32(&client.index, 1)%uint32(len(client.addrs))]
	return
}

// bind
func (client *HTTPClient) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	for _, setOpt := range opts {
		if !setOpt(client) {
			unused = append(unused, setOpt)
		}
	}
	return
}

// update
func (client *HTTPClient) update(addrs []lamp.Address) (err error) {
	client.addrs = addrs
	return
}
