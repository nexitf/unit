package service

import (
	"context"

	"github.com/nexitf/lamp"
	"github.com/nexitf/unit/internal/core/plugin"
)

var (
	id string
)

func init() {
	// Register plugin.
	id = plugin.Register(
		&Plugin{updaters: make(map[string]*Updater)},
	)
}

type Service interface {
	PluginID() string
	update(addrs []lamp.Address) (err error)
}

type Type struct {
}

// PluginID implements plugin.Type.
func (*Type) PluginID() string {
	return id
}

type Updater struct {
	opt    Service
	cancel func()
}

// Bind implements plugin.Plugin.
func (up *Updater) Bind(varp any) {
	up.opt = varp.(Service)
}

// update
func (up *Updater) update(addrs []lamp.Address) (err error) {
	return up.opt.update(addrs)
}

// stop
func (up *Updater) stop() {
	if up.cancel != nil {
		up.cancel()
	}
}

type Plugin struct {
	updaters map[string]*Updater
	lc       *lamp.Client
}

// Name implements plugin.Plugin.
func (plug *Plugin) Name() string {
	return "NexITF service plugin"
}

// Run implements plugin.Plugin.
func (plug *Plugin) Run(ctx context.Context) (err error) {
	plug.lc, err = lamp.NewClient("etcd://127.0.0.1:2379/services")
	return
}

// Stop implements plugin.Plugin.
func (plug *Plugin) Stop(ctx context.Context) (err error) {
	for _, up := range plug.updaters {
		up.stop()
	}
	if plug.lc != nil {
		_ = plug.lc.Close()
	}
	return
}

// Watch implements plugin.Plugin.
func (plug *Plugin) Watch(ctx context.Context, name string, updater plugin.Updater) (err error) {
	up, ok := updater.(*Updater)
	if !ok {
		return
	}
	// Watch
	cancel, err := plug.lc.Watch(name, func(addrs []lamp.Address, closed bool) {
		if !closed {
			up.update(addrs)
		} else {
			up.update(make([]lamp.Address, 0))
		}
	})
	if err != nil {
		return err
	}
	up.cancel = cancel
	plug.updaters[name] = up
	return
}

// NewPlugin implements plugin.Plugin.
func (plug *Plugin) NewUpdater() (updater plugin.Updater) {
	return new(Updater)
}
