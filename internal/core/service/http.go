package service

import (
	"fmt"
	"sync"

	"github.com/nexitf/lamp"
)

type HTTPClient struct {
	Type
	mutex sync.RWMutex
	addrs []lamp.Address
}

func (client *HTTPClient) Get() {

}

func (client *HTTPClient) update(addrs []lamp.Address) (err error) {
	client.addrs = addrs
	fmt.Printf("addrs: %+v\n", addrs)
	return
}
