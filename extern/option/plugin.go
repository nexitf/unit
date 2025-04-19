package option

import (
	"context"
	"errors"

	"github.com/nexitf/unit/extern/plugin"
)

var (
	ErrUnrecognizedBindOption = errors.New("unrecognized bind option")
)

var (
	id   string
	plug *optionPlugin
)

func init() {
	plug = &optionPlugin{updaters: make(map[string]*optionUpdater)}
	// Register plugin.
	id = plugin.Register(plug)
}

type optionBase struct {
}

// UpdaterID implements plugin.Type.
func (base *optionBase) PluginID() string {
	return id
}

// bind implements Option.
func (base *optionBase) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	return
}

type option interface {
	PluginID() string
	bind(opts ...plugin.BindOption) (unused []plugin.BindOption)
	update(value string) (err error)
}

type optionUpdater struct {
	optionBase
	varp option
}

// Bind implements plugin.Updater.
func (up *optionUpdater) Bind(varp any, opts ...plugin.BindOption) {
	up.varp = varp.(option)
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
func (up *optionUpdater) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	for _, setOpt := range opts {
		if !setOpt(up) {
			unused = append(unused, setOpt)
		}
	}
	return
}

// update
func (up *optionUpdater) update(value string) (err error) {
	return up.varp.update(value)
}

type optionPlugin struct {
	updaters map[string]*optionUpdater
}

// Name implements plugin.Plugin.
func (plug *optionPlugin) Name() string {
	return "NexITF option plugin"
}

// Run implements plugin.Plugin.
func (plug *optionPlugin) Run(ctx context.Context) (err error) {
	return
}

// Stop implements plugin.Plugin.
func (plug *optionPlugin) Stop(ctx context.Context) (err error) {
	return
}

// Bind implements plugin.Plugin.
func (plug *optionPlugin) Bind(ctx context.Context, name string, updater plugin.Updater) (err error) {
	up, ok := updater.(*optionUpdater)
	if !ok {
		return
	}
	plug.updaters[name] = up
	return
}

// NewUpdater implements plugin.Plugin.
func (plug *optionPlugin) NewUpdater() (updater plugin.Updater) {
	return new(optionUpdater)
}
