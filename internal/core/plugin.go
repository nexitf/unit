package core

import (
	"github.com/nexitf/unit/internal/core/plugin"
)

// // RunPlugin
// func (u *unit) RunPlugin(ctx context.Context) (err error) {
// 	for _, op := range u.plugins {
// 		if err = op.Run(ctx); err != nil {
// 			return
// 		}
// 	}
// 	return
// }

// // StopPlugin
// func (u *unit) StopPlugin(ctx context.Context) (err error) {
// 	for _, op := range u.plugins {
// 		if err = op.Stop(ctx); err != nil {
// 			return
// 		}
// 	}
// 	return
// }

// LoadPlugin
func (u *unit) LoadPlugin(pluginID string) (plugin.Plugin, bool) {
	if u.plugins == nil {
		u.plugins = plugin.LoadPlugins()
	}
	plugin, exist := u.plugins[pluginID]
	return plugin, exist
}
