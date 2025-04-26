package service

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/nexitf/unit/discovery"
	"github.com/nexitf/unit/internal/errors"
	"github.com/nexitf/unit/plugin"
)

var (
	ErrPluginNotInited          = errors.New("plugin not inited, see 'github.com/nexitf/unit/extern/service.Init()'")
	ErrInvalidUpdater           = errors.New("invalid updater")
	ErrUnrecognizedVariableType = errors.New("unrecognized variable type")
	ErrUnrecognizedBindOption   = errors.New("unrecognized bind option")
)

type Endpoint = discovery.Endpoint

var (
	id   string
	plug *servicePlugin
)

func init() {
	plug = &servicePlugin{updaters: make(map[string]*serviceUpdater)}
	// Register plugin.
	id = plugin.Register(plug)
}

// Resource Type Signature Interface Implementation
type Resource struct {
}

// PluginID implements plugin.Resource.
func (res *Resource) PluginID() string {
	return id
}

type Base struct {
}

// Init implements Service.
func (base *Base) Init() {
}

// Bind implements Service.
func (base *Base) Bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	return opts
}

type Service interface {
	PluginID() string
	Init()
	Bind(opts ...plugin.BindOption) (unused []plugin.BindOption)
	Update(endpoints []Endpoint) (err error)
}

type service interface {
	PluginID() string
	init()
	bind(opts ...plugin.BindOption) (unused []plugin.BindOption)
	update(endpoints []Endpoint) (err error)
}

// WithVariableReady binds a ready function, and supports all service.
func WithVariableReady(fn func()) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		up, used := varp.(*serviceUpdater)
		if used {
			up.readyFn = fn
		}
		return
	}
}

// WithVariableChange binds a change function, and supports all service.
func WithVariableChange(fn func()) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		up, used := varp.(*serviceUpdater)
		if used {
			up.changeFn = fn
		}
		return
	}
}

// WithBeforeUpdate attaches a callback to be invoked before the update.
func WithBeforeUpdate(fn func(endpoints []Endpoint) []Endpoint) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		up, used := varp.(*serviceUpdater)
		if used {
			up.beforeUpdateFn = fn
		}
		return
	}
}

// WithAfterUpdate attaches a callback to be invoked after the update.
func WithAfterUpdate(fn func(endpoints []Endpoint, err error)) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		up, used := varp.(*serviceUpdater)
		if used {
			up.afterUpdateFn = fn
		}
		return
	}
}

// WithTag specifies a tag for filtering the service's endpoints.
func WithTag(tag string) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		up, used := varp.(*serviceUpdater)
		if used && tag != "" {
			up.tag = tag
		}
		return
	}
}

// WithDirectAddr binds the client with the direct addresses.
//
// Support:
//   - HTTPClient
func WithDirectAddr(addrs ...string) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		up, used := varp.(*serviceUpdater)
		if used {
			for _, addr := range addrs {
				up.static = append(up.static, Endpoint{
					ID:     len(up.static) + 1,
					Addr:   addr,
					Weight: 100,
				})
			}
		}
		return
	}
}

// WithDirectEndpoint binds the client with the direct endpoints.
//
// Support:
//   - HTTPClient
func WithDirectEndpoint(endpoints ...Endpoint) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		up, used := varp.(*serviceUpdater)
		if used {
			up.static = append(up.static, endpoints...)
		}
		return
	}
}

type serviceUpdater struct {
	Resource
	mutex          sync.RWMutex
	varp           plugin.Resource
	tag            string
	static         []Endpoint
	uptime         time.Time
	endpoints      []Endpoint
	once           sync.Once
	initFn         func()
	bindFn         func(opts ...plugin.BindOption) (unused []plugin.BindOption)
	updateFn       func(endpoints []Endpoint) (err error)
	beforeUpdateFn func(endpoints []Endpoint) []Endpoint
	afterUpdateFn  func(endpoints []Endpoint, err error)
	readyFn        func()
	changeFn       func()
	cancelFn       func()
}

// Bind implements plugin.Updater.
func (up *serviceUpdater) Bind(varp plugin.Resource, opts ...plugin.BindOption) {
	if varp.PluginID() != id {
		panic(ErrUnrecognizedVariableType)
	}
	// Init updater
	switch vp := varp.(type) {
	case Service:
		up.initFn = vp.Init
		up.bindFn = vp.Bind
		up.updateFn = vp.Update
	case service:
		up.initFn = vp.init
		up.bindFn = vp.bind
		up.updateFn = vp.update
	default:
		panic(ErrUnrecognizedVariableType)
	}
	up.tag = "default"
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
func (up *serviceUpdater) Snapshot() (snapshot plugin.Snapshot) {
	up.mutex.RLock()
	defer up.mutex.RUnlock()
	// Load snapshot
	snapshot.Time = up.uptime
	if up.endpoints == nil {
		snapshot.Data = "[]"
	} else {
		buf, err := json.Marshal(up.endpoints)
		if err != nil {
			snapshot.Data = "[]"
		} else {
			snapshot.Data = string(buf)
		}
	}
	return
}

// bind
func (up *serviceUpdater) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	for _, setOpt := range opts {
		if !setOpt(up) {
			unused = append(unused, setOpt)
		}
	}
	return
}

// update
func (up *serviceUpdater) update(endpoints []Endpoint) (err error) {
	up.mutex.Lock()
	defer up.mutex.Unlock()
	// Before update
	if up.beforeUpdateFn != nil {
		endpoints = up.beforeUpdateFn(endpoints)
	}
	// Update
	err = up.updateFn(endpoints)
	// After update
	if up.afterUpdateFn != nil {
		up.afterUpdateFn(endpoints, err)
	}
	// Update success
	if err == nil {
		up.uptime = time.Now()
		up.endpoints = endpoints
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
func (up *serviceUpdater) stop() {
	if up.cancelFn != nil {
		up.cancelFn()
	}
}

type servicePlugin struct {
	discovery ServiceDiscovery
	updaters  map[string]*serviceUpdater
}

// Name implements plugin.Plugin.
func (plug *servicePlugin) Name() string {
	return "NexITF service plugin"
}

// About implements plugin.Plugin.
func (plug *servicePlugin) About() (about plugin.About) {
	about.Name = plug.Name()
	about.Version = "1.0.0"
	about.Package = "github.com/nexitf/unit/extern/service"
	return
}

// Init
func (plug *servicePlugin) Init(discovery ServiceDiscovery) {
	plug.discovery = discovery
}

// Init inits the plugin.
func Init(discovery ServiceDiscovery) {
	plug.Init(discovery)
}

// Run implements plugin.Plugin.
func (plug *servicePlugin) Run(ctx context.Context) (err error) {
	if plug.discovery == nil {
		return ErrPluginNotInited
	}
	return
}

// Stop implements plugin.Plugin.
func (plug *servicePlugin) Stop(ctx context.Context) (err error) {
	if plug.discovery == nil {
		return ErrPluginNotInited
	}
	for _, up := range plug.updaters {
		up.stop()
	}
	return
}

// Bind implements plugin.Plugin.
func (plug *servicePlugin) Bind(ctx context.Context, name string, updater plugin.Updater) (err error) {
	up, ok := updater.(*serviceUpdater)
	if !ok {
		return ErrInvalidUpdater
	}
	if len(up.static) <= 0 {
		// Watch
		cancel, err := plug.discovery.Watch(ctx, name, up.tag, func(endpoints []Endpoint, _ bool) {
			up.update(endpoints)
		})
		if err != nil {
			return err
		}
		up.cancelFn = cancel
	} else {
		up.update(up.static)
	}
	plug.updaters[name] = up
	return
}

// NewUpdater implements plugin.Plugin.
func (plug *servicePlugin) NewUpdater() (updater plugin.Updater) {
	return new(serviceUpdater)
}

// ServiceDiscovery
type ServiceDiscovery interface {
	Watch(ctx context.Context, serviceName string, tag string, update func(endpoints []discovery.Endpoint, closed bool)) (close func(), err error)
}
