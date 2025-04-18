package option

import (
	"strconv"
	"sync"
)

type String struct {
	Type
	mutex sync.RWMutex
	value string
}

// Get returns the value of the String.
func (s *String) Get() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.value
}

// update updates the value of the String.
func (s *String) update(value string) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.value = value
	return
}

type Int32 struct {
	Type
	mutex sync.RWMutex
	value int32
}

// Get returns the value of the Int32.
func (i *Int32) Get() int32 {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	return i.value
}

// update updates the value of the Int32.
func (i *Int32) update(value string) (err error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i32, err := strconv.ParseInt(value, 10, 32)
	if err == nil {
		i.value = int32(i32)
	}
	return
}

type Int = Int32

type Uint32 struct {
	Type
	mutex sync.RWMutex
	value uint32
}

// Get returns the value of the Uint32.
func (u *Uint32) Get() uint32 {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return u.value
}

// update updates the value of the Uint32.
func (u *Uint32) update(value string) (err error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u32, err := strconv.ParseUint(value, 10, 32)
	if err == nil {
		u.value = uint32(u32)
	}
	return
}

type Int64 struct {
	Type
	mutex sync.RWMutex
	value int64
}

// Get returns the value of the Int64.
func (i *Int64) Get() int64 {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	return i.value
}

// update updates the value of the Int64.
func (i *Int64) update(value string) (err error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.value, err = strconv.ParseInt(value, 10, 64)
	return
}

type Uint64 struct {
	Type
	mutex sync.RWMutex
	value uint64
}

// Get returns the value of the Uint64.
func (u *Uint64) Get() uint64 {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return u.value
}

// update updates the value of the Uint64.
func (u *Uint64) update(value string) (err error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.value, err = strconv.ParseUint(value, 10, 64)
	return
}
