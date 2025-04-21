package service

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/nexitf/lamp"
	"github.com/nexitf/unit/plugin"
)

var (
	ErrPluginNotInited          = errors.New("plugin not inited, see 'github.com/nexitf/unit/extern/service.Init()'")
	ErrInvalidUpdater           = errors.New("invalid updater")
	ErrUnrecognizedVariableType = errors.New("unrecognized variable type")
	ErrUnrecognizedBindOption   = errors.New("unrecognized bind option")
)

var (
	id   string
	plug *servicePlugin
)

func init() {
	plug = &servicePlugin{updaters: make(map[string]*serviceUpdater)}
	// Register plugin.
	id = plugin.Register(plug)
}

type ServiceBase struct {
}

// PluginID implements plugin.Resource.
func (base *ServiceBase) PluginID() string {
	return id
}

// Bind implements Service.
func (base *ServiceBase) Bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	return opts
}

type Service interface {
	PluginID() string
	Bind(opts ...plugin.BindOption) (unused []plugin.BindOption)
	Update(addrs []lamp.Address) (err error)
}

// Service base methods
type serviceBase struct {
}

// PluginID implements plugin.Type.
func (base *serviceBase) PluginID() string {
	return id
}

// bind implements Service.
func (base *serviceBase) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	return opts
}

type service interface {
	PluginID() string
	bind(opts ...plugin.BindOption) (unused []plugin.BindOption)
	update(addrs []lamp.Address) (err error)
}

// WithVariableReady
func WithVariableReady(fn func()) plugin.BindOption {
	return func(v plugin.Resource) (used bool) {
		up, used := v.(*serviceUpdater)
		if used {
			up.readyFn = fn
		}
		return
	}
}

// WithVariableChange
func WithVariableChange(fn func()) plugin.BindOption {
	return func(v plugin.Resource) (used bool) {
		up, used := v.(*serviceUpdater)
		if used {
			up.changeFn = fn
		}
		return
	}
}

type serviceUpdater struct {
	serviceBase
	mutex    sync.RWMutex
	addrs    []lamp.Address
	uptime   time.Time
	varp     plugin.Resource
	once     sync.Once
	bindFn   func(opts ...plugin.BindOption) (unused []plugin.BindOption)
	updateFn func(addrs []lamp.Address) (err error)
	readyFn  func()
	changeFn func()
	cancelFn func()
}

// Bind implements plugin.Updater.
func (up *serviceUpdater) Bind(varp plugin.Resource, opts ...plugin.BindOption) {
	if varp.PluginID() != id {
		panic(ErrUnrecognizedVariableType)
	}
	switch vp := varp.(type) {
	case service:
		up.bindFn = vp.bind
		up.updateFn = vp.update
	case Service:
		up.bindFn = vp.Bind
		up.updateFn = vp.Update
	default:
		panic(ErrUnrecognizedVariableType)
	}
	// Bind options
	unused := up.bind(opts...)
	if len(unused) > 0 {
		unused = up.bindFn(unused...)
	}
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
	if up.addrs == nil {
		snapshot.Data = "[]"
	} else {
		buf, err := json.Marshal(up.addrs)
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
func (up *serviceUpdater) update(addrs []lamp.Address) (err error) {
	up.mutex.Lock()
	defer up.mutex.Unlock()
	// Update
	err = up.updateFn(addrs)
	if err == nil {
		up.addrs = addrs
		up.uptime = time.Now()
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
	updaters map[string]*serviceUpdater
	client   *lamp.Client
}

// Name implements plugin.Plugin.
func (plug *servicePlugin) Name() string {
	return "NexITF service plugin"
}

// About implements plugin.Plugin.
func (plug *servicePlugin) About() (about plugin.About) {
	about.Name = plug.Name()
	about.Version = "1.0.0"
	about.Author = "Kami"
	about.Package = "github.com/nexitf/unit/extern/service"
	return
}

// Init
func (plug *servicePlugin) Init(client *lamp.Client) {
	plug.client = client
}

// Init inits the plugin.
func Init(lc *lamp.Client) {
	plug.Init(lc)
}

// Run implements plugin.Plugin.
func (plug *servicePlugin) Run(ctx context.Context) (err error) {
	if plug.client == nil {
		return ErrPluginNotInited
	}
	return
}

// Stop implements plugin.Plugin.
func (plug *servicePlugin) Stop(ctx context.Context) (err error) {
	if plug.client == nil {
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
	// Watch
	cancel, err := plug.client.Watch(name, func(addrs []lamp.Address, _ bool) {
		up.update(addrs)
	})
	if err != nil {
		return err
	}
	up.cancelFn = cancel
	plug.updaters[name] = up
	return
}

// NewUpdater implements plugin.Plugin.
func (plug *servicePlugin) NewUpdater() (updater plugin.Updater) {
	return new(serviceUpdater)
}
