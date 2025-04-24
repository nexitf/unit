package config

import (
	"strconv"
	"sync"

	"github.com/nexitf/unit/plugin"
)

// WithDefaultString sets the default value of the string option.
func WithDefaultString(value string) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		opt, used := varp.(*String)
		if used {
			opt.value = value
		}
		return
	}
}

type String struct {
	Resource
	mutex sync.RWMutex
	value string
}

// Get returns the value of the String.
func (s *String) Get() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.value
}

// init implements config.
func (s *String) init() {

}

// bind implements config.
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
	return func(varp plugin.Resource) (used bool) {
		opt, used := varp.(*Int)
		if used {
			opt.value = value
		}
		return
	}
}

type Int struct {
	Resource
	mutex sync.RWMutex
	value int
}

// Get returns the value of the Int.
func (i *Int) Get() int {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	return i.value
}

// init implements config.
func (i *Int) init() {

}

// bind implements config.
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
	if value == "" {
		i.value = 0
		return
	}
	iv, err := strconv.Atoi(value)
	if err == nil {
		i.value = iv
	}
	return
}

// WithDefaultInt32 sets the default value of the int32 option.
func WithDefaultInt32(value int32) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		opt, used := varp.(*Int32)
		if used {
			opt.value = value
		}
		return
	}
}

type Int32 struct {
	Resource
	mutex sync.RWMutex
	value int32
}

// Get returns the value of the Int32.
func (i *Int32) Get() int32 {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	return i.value
}

// init implements config.
func (i *Int32) init() {

}

// bind implements config.
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
	if value == "" {
		i.value = 0
		return
	}
	i32, err := strconv.ParseInt(value, 10, 32)
	if err == nil {
		i.value = int32(i32)
	}
	return
}

// WithDefaultUint32 sets the default value of the uint32 option.
func WithDefaultUint32(value uint32) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		opt, used := varp.(*Uint32)
		if used {
			opt.value = value
		}
		return
	}
}

type Uint32 struct {
	Resource
	mutex sync.RWMutex
	value uint32
}

// Get returns the value of the Uint32.
func (u *Uint32) Get() uint32 {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return u.value
}

// init implements config.
func (u *Uint32) init() {

}

// bind implements config.
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
	if value == "" {
		u.value = 0
		return
	}
	u32, err := strconv.ParseUint(value, 10, 32)
	if err == nil {
		u.value = uint32(u32)
	}
	return
}

// WithDefaultInt64 sets the default value of the int64 option.
func WithDefaultInt64(value int64) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		opt, used := varp.(*Int64)
		if used {
			opt.value = value
		}
		return
	}
}

type Int64 struct {
	Resource
	mutex sync.RWMutex
	value int64
}

// Get returns the value of the Int64.
func (i *Int64) Get() int64 {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	return i.value
}

// init implements config.
func (u *Int64) init() {

}

// bind implements config.
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
	if value == "" {
		i.value = 0
		return
	}
	iv, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		i.value = iv
	}
	return
}

// WithDefaultUint64 sets the default value of the uint64 option.
func WithDefaultUint64(value uint64) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		opt, used := varp.(*Uint64)
		if used {
			opt.value = value
		}
		return
	}
}

type Uint64 struct {
	Resource
	mutex sync.RWMutex
	value uint64
}

// Get returns the value of the Uint64.
func (u *Uint64) Get() uint64 {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return u.value
}

// init implements config.
func (u *Uint64) init() {

}

// bind implements config.
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
	if value == "" {
		u.value = 0
		return
	}
	uv, err := strconv.ParseUint(value, 10, 64)
	if err == nil {
		u.value = uv
	}
	return
}

// WithDefaultFloat32 sets the default value of the float32 option.
func WithDefaultFloat32(value float32) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		opt, used := varp.(*Float32)
		if used {
			opt.value = value
		}
		return
	}
}

type Float32 struct {
	Resource
	mutex sync.RWMutex
	value float32
}

// Get returns the value of the Float32.
func (f *Float32) Get() float32 {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.value
}

// init implements config.
func (f *Float32) init() {

}

// bind implements config.
func (f *Float32) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	for _, setOpt := range opts {
		if !setOpt(f) {
			unused = append(unused, setOpt)
		}
	}
	return
}

// update updates the value of the Float32.
func (f *Float32) update(value string) (err error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	if value == "" {
		f.value = 0
		return
	}
	f32, err := strconv.ParseFloat(value, 64)
	if err == nil {
		f.value = float32(f32)
	}
	return
}

// WithDefaultFloat64 sets the default value of the float64 option.
func WithDefaultFloat64(value float64) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		opt, used := varp.(*Float64)
		if used {
			opt.value = value
		}
		return
	}
}

type Float64 struct {
	Resource
	mutex sync.RWMutex
	value float64
}

// Get returns the value of the Float64.
func (f *Float64) Get() float64 {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.value
}

// init implements config.
func (f *Float64) init() {

}

// bind implements config.
func (f *Float64) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	for _, setOpt := range opts {
		if !setOpt(f) {
			unused = append(unused, setOpt)
		}
	}
	return
}

// update updates the value of the Float64.
func (f *Float64) update(value string) (err error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	if value == "" {
		f.value = 0
		return
	}
	f64, err := strconv.ParseFloat(value, 64)
	if err == nil {
		f.value = f64
	}
	return
}
