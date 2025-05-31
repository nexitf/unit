package service

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/nexitf/logkit"
	"github.com/nexitf/unit/analytics/stats"
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
	plug = &servicePlugin{
		discovery: discovery.Offline,
		updaters:  make(map[string]*serviceUpdater),
	}
}

// PluginID returns the plugin ID.
func PluginID() string {
	return id
}

// Resource Type Signature Interface Implementation
type Resource struct {
}

// PluginID implements plugin.Resource.
func (res *Resource) PluginID() (pluginID string, pluginName string) {
	return id, plug.Name()
}

type Base struct {
	Resource
}

// Init implements Client.
func (base *Base) Init() {
}

// Bind implements Client.
func (base *Base) Bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	return opts
}

// Update implements Client.
func (base *Base) Update(endpoints []Endpoint) (err error) {
	return
}

// Close implements Client.
func (base *Base) Close() (err error) {
	return
}

type Client interface {
	PluginID() (pluginID string, pluginName string)
	Init()
	Bind(opts ...plugin.BindOption) (unused []plugin.BindOption)
	Update(endpoints []Endpoint) (err error)
	Close() (err error)
}

type client interface {
	PluginID() (pluginID string, pluginName string)
	init()
	bind(opts ...plugin.BindOption) (unused []plugin.BindOption)
	update(endpoints []Endpoint) (err error)
	close() (err error)
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
	closeFn        func() (err error)
	readyFn        func()
	changeFn       func()
	cancelFn       func()
}

// Bind implements plugin.Updater.
func (up *serviceUpdater) Bind(varp plugin.Resource, opts ...plugin.BindOption) {
	pluginID, _ := varp.PluginID()
	if pluginID != id {
		panic(ErrUnrecognizedVariableType)
	}
	// Init updater
	switch vp := varp.(type) {
	case Client:
		up.initFn = vp.Init
		up.bindFn = vp.Bind
		up.updateFn = vp.Update
		up.closeFn = vp.Close
	case client:
		up.initFn = vp.init
		up.bindFn = vp.bind
		up.updateFn = vp.update
		up.closeFn = vp.close
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
	if up.closeFn != nil {
		spend := stats.CallSpend(func() {
			up.closeFn()
		})
		logkit.Info("client close", logkit.Field("spend", spend.Seconds()))
	}
	if up.cancelFn != nil {
		spend := stats.CallSpend(func() {
			up.cancelFn()
		})
		logkit.Info("watcher cancel", logkit.Field("spend", spend.Seconds()))
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

// Plugin initializes and returns a service plugin instance.
// It takes an implementation of the ServiceDiscovery interface as a parameter
// and initializes the global service plugin instance.
// Parameters:
//   - discovery: An implementation of the ServiceDiscovery interface for service discovery operations.
//
// Returns:
//   - plugin.Plugin: An initialized service plugin instance implementing the plugin.Plugin interface.
func Plugin(discovery ServiceDiscovery) plugin.Plugin {
	plug.Init(discovery)
	return plug
}

// Run implements plugin.Plugin.
func (plug *servicePlugin) Run(ctx context.Context, pluginID string) (err error) {
	if plug.discovery == nil {
		return ErrPluginNotInited
	}
	id = pluginID
	return
}

// Stop implements plugin.Plugin.
func (plug *servicePlugin) Stop(ctx context.Context) (err error) {
	if plug.discovery == nil {
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
func (plug *servicePlugin) Bind(ctx context.Context, name string, updater plugin.Updater) (err error) {
	up, ok := updater.(*serviceUpdater)
	if !ok {
		return ErrInvalidUpdater
	}
	if len(up.static) <= 0 {
		// Watch
		cancel, err := plug.discovery.Watch(ctx, name, up.tag, func(endpoints []Endpoint, closed bool) {
			if len(endpoints) <= 0 {
				if closed {
					up.update(nil)
				}
			} else {
				up.update(endpoints)
			}
		})
		if err != nil {
			return err
		}
		up.cancelFn = func() { cancel() }
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
	Watch(ctx context.Context, serviceName, tag string, update func(endpoints []discovery.Endpoint, closed bool)) (cancel func() error, err error)
}
