package service

import (
	"context"
	"errors"

	"github.com/nexitf/lamp"
	"github.com/nexitf/unit/extern/plugin"
)

var (
	ErrPluginNotInited        = errors.New("plugin not inited, see 'github.com/nexitf/unit/extern/service.Init()'")
	ErrInvalidUpdater         = errors.New("invalid updater")
	ErrUnrecognizedBindOption = errors.New("unrecognized bind option")
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

// Service base methods
type serviceBase struct {
}

// PluginID implements plugin.Type.
func (base *serviceBase) PluginID() string {
	return id
}

// bind implements Service.
func (base *serviceBase) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	return
}

type service interface {
	PluginID() string
	bind(opts ...plugin.BindOption) (unused []plugin.BindOption)
	update(addrs []lamp.Address) (err error)
}

type serviceUpdater struct {
	serviceBase
	varp   service
	cancel func()
}

// Bind implements plugin.Plugin.
func (up *serviceUpdater) Bind(varp any, opts ...plugin.BindOption) {
	up.varp = varp.(service)
	// Bind options
	unused := up.varp.bind(opts...)
	if len(unused) > 0 {
		unused = up.bind(unused...)
	}
	if len(unused) > 0 {
		panic(ErrUnrecognizedBindOption)
	}
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
	return up.varp.update(addrs)
}

// stop
func (up *serviceUpdater) stop() {
	if up.cancel != nil {
		up.cancel()
	}
}

type servicePlugin struct {
	updaters map[string]*serviceUpdater
	lc       *lamp.Client
}

// Name implements plugin.Plugin.
func (plug *servicePlugin) Name() string {
	return "NexITF service plugin"
}

// Init
func (plug *servicePlugin) Init(lc *lamp.Client) {
	plug.lc = lc
}

// Init inits the plugin.
func Init(lc *lamp.Client) {
	plug.Init(lc)
}

// Run implements plugin.Plugin.
func (plug *servicePlugin) Run(ctx context.Context) (err error) {
	if plug.lc == nil {
		return ErrPluginNotInited
	}
	return
}

// Stop implements plugin.Plugin.
func (plug *servicePlugin) Stop(ctx context.Context) (err error) {
	if plug.lc == nil {
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
	cancel, err := plug.lc.Watch(name, func(addrs []lamp.Address, _ bool) {
		up.update(addrs)
	})
	if err != nil {
		return err
	}
	up.cancel = cancel
	plug.updaters[name] = up
	return
}

// NewUpdater implements plugin.Plugin.
func (plug *servicePlugin) NewUpdater() (updater plugin.Updater) {
	return new(serviceUpdater)
}
