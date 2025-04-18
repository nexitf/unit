package option

import (
	"context"

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

type Option interface {
	PluginID() string
	update(value string) (err error)
}

type Type struct {
}

// UpdaterID implements plugin.Type.
func (*Type) PluginID() string {
	return id
}

type Updater struct {
	opt Option
}

// Bind implements plugin.Updater.
func (up *Updater) Bind(varp any) {
	up.opt = varp.(Option)
}

// update
func (up *Updater) update(value string) (err error) {
	return up.opt.update(value)
}

type Plugin struct {
	updaters map[string]*Updater
}

// Name implements plugin.Plugin.
func (plug *Plugin) Name() string {
	return "NexITF option plugin"
}

// Run implements plugin.Plugin.
func (plug *Plugin) Run(ctx context.Context) (err error) {
	return
}

// Stop implements plugin.Plugin.
func (plug *Plugin) Stop(ctx context.Context) (err error) {
	return
}

// Watch implements plugin.Plugin.
func (plug *Plugin) Watch(ctx context.Context, name string, updater plugin.Updater) (err error) {
	up, ok := updater.(*Updater)
	if !ok {
		return
	}
	plug.updaters[name] = up
	return
}

// NewUpdater implements plugin.Plugin.
func (plug *Plugin) NewUpdater() (updater plugin.Updater) {
	return new(Updater)
}
