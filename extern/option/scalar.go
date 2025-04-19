package option

import (
	"strconv"
	"sync"

	"github.com/nexitf/unit/extern/plugin"
)

// WithDefaultString sets the default value of the string option.
func WithDefaultString(value string) plugin.BindOption {
	return func(v plugin.Variable) (used bool) {
		opt, used := v.(*String)
		if used {
			opt.value = value
		}
		return
	}
}

type String struct {
	optionBase
	mutex sync.RWMutex
	value string
}

// Get returns the value of the String.
func (s *String) Get() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.value
}

// bind
func (s *String) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	for _, setOpt := range opts {
		if !setOpt(s) {
			unused = append(unused, setOpt)
		}
	}
	return
}

// update updates the value of the String.
func (s *String) update(value string) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.value = value
	return
}

// WithDefaultInt sets the default value of the int option.
func WithDefaultInt(value int) plugin.BindOption {
	return func(v plugin.Variable) (used bool) {
		opt, used := v.(*Int)
		if used {
			opt.value = value
		}
		return
	}
}

type Int struct {
	optionBase
	mutex sync.RWMutex
	value int
}

// Get returns the value of the Int.
func (i *Int) Get() int {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	return i.value
}

// bind
func (i *Int) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	for _, setOpt := range opts {
		if !setOpt(i) {
			unused = append(unused, setOpt)
		}
	}
	return
}

// update updates the value of the Int.
func (i *Int) update(value string) (err error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	iv, err := strconv.Atoi(value)
	if err == nil {
		i.value = iv
	}
	return
}

// WithDefaultInt32 sets the default value of the int32 option.
func WithDefaultInt32(value int32) plugin.BindOption {
	return func(v plugin.Variable) (used bool) {
		opt, used := v.(*Int32)
		if used {
			opt.value = value
		}
		return
	}
}

type Int32 struct {
	optionBase
	mutex sync.RWMutex
	value int32
}

// Get returns the value of the Int32.
func (i *Int32) Get() int32 {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	return i.value
}

// bind
func (i *Int32) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	for _, setOpt := range opts {
		if !setOpt(i) {
			unused = append(unused, setOpt)
		}
	}
	return
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

// WithDefaultUint32 sets the default value of the uint32 option.
func WithDefaultUint32(value uint32) plugin.BindOption {
	return func(v plugin.Variable) (used bool) {
		opt, used := v.(*Uint32)
		if used {
			opt.value = value
		}
		return
	}
}

type Uint32 struct {
	optionBase
	mutex sync.RWMutex
	value uint32
}

// Get returns the value of the Uint32.
func (u *Uint32) Get() uint32 {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return u.value
}

// bind
func (u *Uint32) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	for _, setOpt := range opts {
		if !setOpt(u) {
			unused = append(unused, setOpt)
		}
	}
	return
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

// WithDefaultInt64 sets the default value of the int64 option.
func WithDefaultInt64(value int64) plugin.BindOption {
	return func(v plugin.Variable) (used bool) {
		opt, used := v.(*Int64)
		if used {
			opt.value = value
		}
		return
	}
}

type Int64 struct {
	optionBase
	mutex sync.RWMutex
	value int64
}

// Get returns the value of the Int64.
func (i *Int64) Get() int64 {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	return i.value
}

// bind
func (i *Int64) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	for _, setOpt := range opts {
		if !setOpt(i) {
			unused = append(unused, setOpt)
		}
	}
	return
}

// update updates the value of the Int64.
func (i *Int64) update(value string) (err error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.value, err = strconv.ParseInt(value, 10, 64)
	return
}

// WithDefaultUint64 sets the default value of the uint64 option.
func WithDefaultUint64(value uint64) plugin.BindOption {
	return func(v plugin.Variable) (used bool) {
		opt, used := v.(*Uint64)
		if used {
			opt.value = value
		}
		return
	}
}

type Uint64 struct {
	optionBase
	mutex sync.RWMutex
	value uint64
}

// Get returns the value of the Uint64.
func (u *Uint64) Get() uint64 {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return u.value
}

// bind
func (u *Uint64) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	for _, setOpt := range opts {
		if !setOpt(u) {
			unused = append(unused, setOpt)
		}
	}
	return
}

// update updates the value of the Uint64.
func (u *Uint64) update(value string) (err error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.value, err = strconv.ParseUint(value, 10, 64)
	return
}
