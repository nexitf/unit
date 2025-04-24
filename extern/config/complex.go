package config

import (
	"strconv"
	"strings"
	"sync"

	"github.com/nexitf/unit/plugin"
)

type Strings struct {
	Resource
	mutex sync.RWMutex
	value []string
}

// Get returns the value of the Strings.
func (s *Strings) Get() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.value
}

// init implements config.
func (s *Strings) init() {

}

// bind implements config.
func (s *Strings) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	for _, setOpt := range opts {
		if !setOpt(s) {
			unused = append(unused, setOpt)
		}
	}
	return
}

// update updates the value of the String.
func (s *Strings) update(value string) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if value == "" {
		s.value = nil
		return
	}
	s.value = strings.Split(value, ",")
	return
}

type Int64s struct {
	Resource
	mutex sync.RWMutex
	value []int64
}

// Get returns the value of the Int64s.
func (i *Int64s) Get() []int64 {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	return i.value
}

// init implements config.
func (i *Int64s) init() {

}

// bind implements config.
func (i *Int64s) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	for _, setOpt := range opts {
		if !setOpt(i) {
			unused = append(unused, setOpt)
		}
	}
	return
}

// update updates the value of the String.
func (i *Int64s) update(value string) (err error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if value == "" {
		i.value = nil
		return
	}
	var is []int64
	var ss = strings.Split(value, ",")
	for _, s := range ss {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		} else {
			is = append(is, v)
		}
	}
	i.value = is
	return
}
