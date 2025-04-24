package grpc

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/nexitf/lamp"
	"github.com/nexitf/unit/internal/errors"
	"github.com/nexitf/unit/plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

var (
	ErrPluginNotInited          = errors.New("plugin not inited, see 'github.com/nexitf/unit/extern/grpc.Init()'")
	ErrInvalidUpdater           = errors.New("invalid updater")
	ErrUnrecognizedVariableType = errors.New("unrecognized variable type")
	ErrUnrecognizedBindOption   = errors.New("unrecognized bind option")
)

type ClientConn = grpc.ClientConn
type Endpoint = lamp.Endpoint

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

type Service interface {
	PluginID() string
	init()
	bind(opts ...plugin.BindOption) (unused []plugin.BindOption)
	dial(scheme, name string, builder resolver.Builder) (cc *grpc.ClientConn, err error)
	Dialed(cc *grpc.ClientConn)
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

// WithDirectAddr binds the client with the direct addresses.
func WithDirectAddr(addr string) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		up, used := varp.(*serviceUpdater)
		if used {
			up.static = append(up.static, Endpoint{ID: 1, Addr: addr, Weight: 100})
		}
		return
	}
}

type serviceUpdater struct {
	Resource
	mutex     sync.RWMutex
	varp      plugin.Resource
	static    []Endpoint
	uptime    time.Time
	endpoints []Endpoint
	once      sync.Once
	initFn    func()
	bindFn    func(opts ...plugin.BindOption) (unused []plugin.BindOption)
	dialFn    func(scheme, name string, builder resolver.Builder) (cc *grpc.ClientConn, err error)
	dialedFn  func(cc *grpc.ClientConn)
	updateFn  func(endpoints []Endpoint) (err error)
	readyFn   func()
	changeFn  func()
	cancelFn  func()
}

// Bind implements plugin.Updater.
func (up *serviceUpdater) Bind(varp plugin.Resource, opts ...plugin.BindOption) {
	if varp.PluginID() != id {
		panic(ErrUnrecognizedVariableType)
	}
	switch vp := varp.(type) {
	case Service:
		up.initFn = vp.init
		up.bindFn = vp.bind
		up.dialFn = vp.dial
		up.dialedFn = vp.Dialed
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

// dial
func (up *serviceUpdater) dial(scheme, name string, builder resolver.Builder) (cc *grpc.ClientConn, err error) {
	return up.dialFn(scheme, name, builder)
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
	// Update
	err = up.updateFn(endpoints)
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

// ResolveNow implements resolver.Resolver.
func (up *serviceUpdater) ResolveNow(resolver.ResolveNowOptions) {

}

// Close implements resolver.Resolver.
func (up *serviceUpdater) Close() {
	up.stop()
}

type servicePlugin struct {
	updaters  map[string]*serviceUpdater
	discovery ServiceDiscovery
}

// Name implements plugin.Plugin.
func (plug *servicePlugin) Name() string {
	return "NexITF grpc plugin"
}

// About implements plugin.Plugin.
func (plug *servicePlugin) About() (about plugin.About) {
	about.Name = plug.Name()
	about.Version = "1.0.0"
	about.Author = "Kami"
	about.Package = "github.com/nexitf/unit/extern/grpc"
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
	var cc *grpc.ClientConn
	// Dial
	if len(up.static) <= 0 {
		cc, err = up.dial("lamp", name, plug)
	} else {
		cc, err = up.dial("passthrough", up.static[0].Addr, plug)
		if err == nil {
			up.updateFn = func(endpoints []Endpoint) (err error) { return }
		}
	}
	if err != nil {
		return
	}
	plug.updaters[name] = up
	//
	up.dialedFn(cc)
	// Update static
	if len(up.static) > 0 {
		up.update(up.static)
	}
	// Ready
	if up.readyFn != nil {
		up.once.Do(up.readyFn)
	}
	return
}

// NewUpdater implements plugin.Plugin.
func (plug *servicePlugin) NewUpdater() (updater plugin.Updater) {
	return new(serviceUpdater)
}

// Scheme implements resolver.Builder.
func (plug *servicePlugin) Scheme() string {
	return "lamp"
}

// Build implements resolver.Builder.
func (plug *servicePlugin) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var (
		name = target.Endpoint()
		up   = plug.updaters[name]
	)
	up.updateFn = func(endpoints []Endpoint) (err error) {
		var addrs []resolver.Address
		for _, endpoint := range endpoints {
			addrs = append(addrs, resolver.Address{Addr: endpoint.Addr})
		}
		return cc.UpdateState(resolver.State{Addresses: addrs})
	}
	// Watch
	cancel, err := plug.discovery.Watch(name, func(endpoints []Endpoint, _ bool) {
		up.update(endpoints)
	})
	if err != nil {
		up.updateFn = nil
		return nil, err
	}
	up.cancelFn = cancel
	return up, nil
}

// ServiceDiscovery
type ServiceDiscovery interface {
	Watch(serviceName string, update func(endpoints []Endpoint, closed bool)) (close func(), err error)
}
