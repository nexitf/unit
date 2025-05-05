package config

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nexitf/logkit"
	"github.com/nexitf/unit/analytics/stats"
	"github.com/nexitf/unit/plugin"
)

var (
	ErrPluginNotInited          = errors.New("plugin not inited, see 'github.com/nexitf/unit/extern/config.Init()'")
	ErrInvalidUpdater           = errors.New("invalid updater")
	ErrUnrecognizedVariableType = errors.New("unrecognized variable type")
	ErrUnrecognizedBindOption   = errors.New("unrecognized bind option")
)

var (
	id   string
	plug *configPlugin
)

func init() {
	plug = &configPlugin{updaters: make(map[string]*configUpdater)}
	// Register plugin.
	id = plugin.Register(plug)
}

// PluginID returns the plugin ID.
func PluginID() string {
	return id
}

// Resource Type Signature Interface Implementation
type Resource struct {
}

// PluginID implements plugin.Resource.
func (res *Resource) PluginID() string {
	return id
}

type Base struct {
	Resource
}

// Init implements Option.
func (base *Base) Init() {

}

// Bind implements Option.
func (base *Base) Bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	return opts
}

type Config interface {
	PluginID() string
	Init()
	Bind(opts ...plugin.BindOption) (unused []plugin.BindOption)
	Update(value string) (err error)
}

type config interface {
	PluginID() string
	init()
	bind(opts ...plugin.BindOption) (unused []plugin.BindOption)
	update(value string) (err error)
}

// WithVariableReady binds a ready function, and supports all config.
func WithVariableReady(fn func()) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		up, used := varp.(*configUpdater)
		if used {
			up.readyFn = fn
		}
		return
	}
}

// WithVariableChange binds a change function, and supports all config.
func WithVariableChange(fn func()) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		up, used := varp.(*configUpdater)
		if used {
			up.changeFn = fn
		}
		return
	}
}

// WithBeforeUpdate attaches a callback to be invoked before the update.
func WithBeforeUpdate(fn func(value string) string) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		up, used := varp.(*configUpdater)
		if used {
			up.beforeUpdateFn = fn
		}
		return
	}
}

// WithAfterUpdate attaches a callback to be invoked after the update.
func WithAfterUpdate(fn func(value string, err error)) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		up, used := varp.(*configUpdater)
		if used {
			up.afterUpdateFn = fn
		}
		return
	}
}

type configUpdater struct {
	Resource
	mutex          sync.RWMutex
	varp           plugin.Resource
	uptime         time.Time
	value          string
	once           sync.Once
	initFn         func()
	bindFn         func(opts ...plugin.BindOption) (unused []plugin.BindOption)
	updateFn       func(value string) (err error)
	beforeUpdateFn func(value string) string
	afterUpdateFn  func(value string, err error)
	readyFn        func()
	changeFn       func()
}

// Bind implements plugin.Updater.
func (up *configUpdater) Bind(varp plugin.Resource, opts ...plugin.BindOption) {
	if varp.PluginID() != id {
		panic(ErrUnrecognizedVariableType)
	}
	// Init updater
	switch vp := varp.(type) {
	case Config:
		up.initFn = vp.Init
		up.bindFn = vp.Bind
		up.updateFn = vp.Update
	case config:
		up.initFn = vp.init
		up.bindFn = vp.bind
		up.updateFn = vp.update
	default:
		panic(ErrUnrecognizedVariableType)
	}
	// Init variable
	up.initFn()
	// Bind options
	unused := up.bindFn(up.bind(opts...)...)
	if len(unused) > 0 {
		panic(ErrUnrecognizedBindOption)
	}
	up.varp = varp
}

// Snapshot implements plugin.Updater.
func (up *configUpdater) Snapshot() (snapshot plugin.Snapshot) {
	up.mutex.RLock()
	defer up.mutex.RUnlock()
	// Load snapshot
	snapshot.Time = up.uptime
	snapshot.Data = up.value
	return
}

// bind
func (up *configUpdater) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	for _, setOpt := range opts {
		if !setOpt(up) {
			unused = append(unused, setOpt)
		}
	}
	return
}

// update
func (up *configUpdater) update(value string) (err error) {
	up.mutex.Lock()
	defer up.mutex.Unlock()
	// Before update
	if up.beforeUpdateFn != nil {
		value = up.beforeUpdateFn(value)
	}
	// Update
	err = up.updateFn(value)
	// After update
	if up.afterUpdateFn != nil {
		up.afterUpdateFn(value, err)
	}
	// Update success
	if err == nil {
		up.uptime = time.Now()
		up.value = value
		if up.readyFn != nil {
			up.once.Do(up.readyFn)
		}
		if up.changeFn != nil {
			up.changeFn()
		}
	}
	return
}

// stop
func (up *configUpdater) stop() {
}

type configPlugin struct {
	updaters map[string]*configUpdater
	storage  Storage
}

// Name implements plugin.Plugin.
func (plug *configPlugin) Name() string {
	return "NexITF config plugin"
}

// About implements plugin.Plugin.
func (plug *configPlugin) About() (about plugin.About) {
	about.Name = plug.Name()
	about.Version = "1.0.0"
	about.Description = fmt.Sprintf("Powered by: %s(%s@%s)", plug.storage.Name(), plug.storage.Package(), plug.storage.Version())
	about.Package = "github.com/nexitf/unit/extern/config"
	return
}

// Init
func (plug *configPlugin) Init(storage Storage) {
	plug.storage = storage
}

// Init inits the plugin.
func Init(storage Storage) {
	plug.Init(storage)
}

// Run implements plugin.Plugin.
func (plug *configPlugin) Run(ctx context.Context) (err error) {
	if plug.storage == nil {
		return ErrPluginNotInited
	}
	return
}

// Stop implements plugin.Plugin.
func (plug *configPlugin) Stop(ctx context.Context) (err error) {
	if plug.storage == nil {
		return ErrPluginNotInited
	}
	for name, up := range plug.updaters {
		spend := stats.CallSpend(func() {
			up.stop()
		})
		logkit.Info("updater stop successfully",
			logkit.Field("name", name),
			logkit.Field("spend", spend.Seconds()),
		)
	}
	return
}

// Bind implements plugin.Plugin.
func (plug *configPlugin) Bind(ctx context.Context, name string, updater plugin.Updater) (err error) {
	up, ok := updater.(*configUpdater)
	if !ok {
		return ErrInvalidUpdater
	}
	// Watch
	err = plug.storage.Watch(name, func(value string) {
		up.update(strings.TrimSpace(value))
	})
	if err != nil {
		return
	}
	plug.updaters[name] = up
	return
}

// NewUpdater implements plugin.Plugin.
func (plug *configPlugin) NewUpdater() (updater plugin.Updater) {
	return new(configUpdater)
}

type Storage interface {
	// Name
	Name() string

	// Version
	Version() string

	// Package
	Package() string

	// Watch
	Watch(name string, update func(value string)) (err error)
}
