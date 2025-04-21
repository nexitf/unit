package option

// import (
// 	"context"
// 	"errors"
// 	"sync"

// 	"github.com/nexitf/unit/extern/plugin"
// )

// var (
// 	ErrUnrecognizedVariableType = errors.New("unrecognized variable type")
// 	ErrUnrecognizedBindOption   = errors.New("unrecognized bind option")
// )

// var (
// 	id   string
// 	plug *optionPlugin
// )

// func init() {
// 	plug = &optionPlugin{updaters: make(map[string]*optionUpdater)}
// 	// Register plugin.
// 	id = plugin.Register(plug)
// }

// type optionBase struct {
// }

// // UpdaterID implements plugin.Type.
// func (base *optionBase) PluginID() string {
// 	return id
// }

// // bind implements Option.
// func (base *optionBase) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
// 	return
// }

// type option interface {
// 	PluginID() string
// 	bind(opts ...plugin.BindOption) (unused []plugin.BindOption)
// 	update(value string) (err error)
// }

// // WithVariableReady
// func WithVariableReady(fn func()) plugin.BindOption {
// 	return func(v plugin.Variable) (used bool) {
// 		up, used := v.(*optionUpdater)
// 		if used {
// 			up.ready = fn
// 		}
// 		return
// 	}
// }

// // WithVariableChange
// func WithVariableChange(fn func()) plugin.BindOption {
// 	return func(v plugin.Variable) (used bool) {
// 		up, used := v.(*optionUpdater)
// 		if used {
// 			up.change = fn
// 		}
// 		return
// 	}
// }

// type optionUpdater struct {
// 	optionBase
// 	varp   option
// 	once   sync.Once
// 	ready  func()
// 	change func()
// }

// // Bind implements plugin.Updater.
// func (up *optionUpdater) Bind(varp any, opts ...plugin.BindOption) {
// 	vp, ok := varp.(option)
// 	if !ok {
// 		panic(ErrUnrecognizedVariableType)
// 	}
// 	// Bind options
// 	unused := vp.bind(opts...)
// 	if len(unused) > 0 {
// 		unused = up.bind(unused...)
// 	}
// 	if len(unused) > 0 {
// 		panic(ErrUnrecognizedBindOption)
// 	}
// 	up.varp = vp
// }

// // bind
// func (up *optionUpdater) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
// 	for _, setOpt := range opts {
// 		if !setOpt(up) {
// 			unused = append(unused, setOpt)
// 		}
// 	}
// 	return
// }

// // update
// func (up *optionUpdater) update(value string) (err error) {
// 	err = up.varp.update(value)
// 	if err == nil {
// 		if up.ready != nil {
// 			up.once.Do(up.ready)
// 		}
// 		if up.change != nil {
// 			up.change()
// 		}
// 	}
// 	return
// }

// type optionPlugin struct {
// 	updaters map[string]*optionUpdater
// }

// // Name implements plugin.Plugin.
// func (plug *optionPlugin) Name() string {
// 	return "NexITF option plugin"
// }

// // About implements plugin.Plugin.
// func (plug *optionPlugin) About() (about plugin.About) {
// 	about.Name = plug.Name()
// 	about.Author = "Kami"
// 	about.Package = "github.com/nexitf/unit/extern/option"
// 	return
// }

// // Run implements plugin.Plugin.
// func (plug *optionPlugin) Run(ctx context.Context) (err error) {
// 	return
// }

// // Stop implements plugin.Plugin.
// func (plug *optionPlugin) Stop(ctx context.Context) (err error) {
// 	return
// }

// // Bind implements plugin.Plugin.
// func (plug *optionPlugin) Bind(ctx context.Context, name string, updater plugin.Updater) (err error) {
// 	up, ok := updater.(*optionUpdater)
// 	if !ok {
// 		return
// 	}
// 	plug.updaters[name] = up
// 	return
// }

// // NewUpdater implements plugin.Plugin.
// func (plug *optionPlugin) NewUpdater() (updater plugin.Updater) {
// 	return new(optionUpdater)
// }
