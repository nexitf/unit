package service

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/nexitf/lamp"
	"github.com/nexitf/unit/internal/errors"
	"github.com/nexitf/unit/plugin"
)

var (
	ErrPluginNotInited          = errors.New("plugin not inited, see 'github.com/nexitf/unit/extern/service.Init()'")
	ErrInvalidUpdater           = errors.New("invalid updater")
	ErrUnrecognizedVariableType = errors.New("unrecognized variable type")
	ErrUnrecognizedBindOption   = errors.New("unrecognized bind option")
)

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

type Base struct {
}

// PluginID implements plugin.Resource.
func (base *Base) PluginID() string {
	return id
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

// Service base methods
type base struct {
}

// PluginID implements plugin.Resource.
func (base *base) PluginID() string {
	return id
}

// init implements service.
func (base *base) init() {
}

// bind implements service.
func (base *base) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	return opts
}

type service interface {
	PluginID() string
	init()
	bind(opts ...plugin.BindOption) (unused []plugin.BindOption)
	update(endpoints []Endpoint) (err error)
}

// WithVariableReady binds a ready function, support all service.
func WithVariableReady(fn func()) plugin.BindOption {
	return func(v plugin.Resource) (used bool) {
		up, used := v.(*serviceUpdater)
		if used {
			up.readyFn = fn
		}
		return
	}
}

// WithVariableChange binds a change function, support all service.
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
	base
	mutex     sync.RWMutex
	varp      plugin.Resource
	endpoints []Endpoint
	uptime    time.Time
	once      sync.Once
	initFn    func()
	bindFn    func(opts ...plugin.BindOption) (unused []plugin.BindOption)
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
	case service:
		up.initFn = vp.init
		up.bindFn = vp.bind
		up.updateFn = vp.update
	case Service:
		up.initFn = vp.Init
		up.bindFn = vp.Bind
		up.updateFn = vp.Update
	default:
		panic(ErrUnrecognizedVariableType)
	}
	// Init variable
	up.initFn()
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
	// Update
	err = up.updateFn(endpoints)
	if err == nil {
		up.endpoints = endpoints
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
	cancel, err := plug.client.Watch(name, func(endpoints []Endpoint, _ bool) {
		up.update(endpoints)
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
