package core

import (
	"context"

	"github.com/nexitf/logkit"
	"github.com/nexitf/unit/internal/core/plugin"
	"github.com/nexitf/unit/internal/core/utils"
)

// RunPlugins sequentially executes the Run method of all plugins in the Unit instance.
// This method iterates over all registered plugins in the Unit, calls the Run method of each plugin in order,
// and records the execution time of each plugin. If an error occurs during the execution of a plugin,
// it logs the error and returns immediately.
//
// Parameters:
// - ctx: The context used to control the lifecycle of plugin execution, which can be used for cancellation or timeout control.
//
// Returns:
// - err: Returns nil if all plugins are executed successfully; otherwise, returns the error.
func (u *Unit) RunPlugins(ctx context.Context) (err error) {
	for pluginID, plugin := range u.plugins {
		spend := utils.CallSpend(func() {
			err = plugin.Run(ctx, pluginID)
		})
		if err != nil {
			logkit.ErrorWrap(err, "plugin run failed",
				logkit.Field("plugin", plugin.Name()),
				logkit.Field("spend", spend.Seconds()),
			)
			return
		}
		logkit.Info("plugin run successfully",
			logkit.Field("plugin", plugin.Name()),
			logkit.Field("spend", spend.Seconds()),
		)
	}
	return
}

// StopPlugins sequentially executes the Stop method of all plugins in the Unit instance.
// This method iterates over all registered plugins in the Unit, calls the Stop method of each plugin in order,
// and records the execution time of each plugin. If an error occurs during the execution of a plugin,
// it logs the error and returns immediately.
//
// Parameters:
// - ctx: The context used to control the lifecycle of plugin stopping, which can be used for cancellation or timeout control.
//
// Returns:
// - err: Returns nil if all plugins are stopped successfully; otherwise, returns the error.
func (u *Unit) StopPlugins(ctx context.Context) (err error) {
	for _, plugin := range u.plugins {
		spend := utils.CallSpend(func() {
			err = plugin.Stop(ctx)
		})
		if err != nil {
			logkit.ErrorWrap(err, "plugin stop failed",
				logkit.Field("plugin", plugin.Name()),
				logkit.Field("spend", spend.Seconds()),
			)
			return
		}
		logkit.Info("plugin stop successfully",
			logkit.Field("plugin", plugin.Name()),
			logkit.Field("spend", spend.Seconds()),
		)
	}
	return
}

// FindPlugin searches for a plugin with the specified ID in the Unit instance.
// It looks up the plugin in the Unit's plugin map using the provided plugin ID.
//
// Parameters:
//   - pluginID: The unique identifier of the plugin to search for.
//
// Returns:
//   - plugin: The plugin instance if found; otherwise, the zero value of plugin.Plugin.
//   - found: A boolean indicating whether the plugin was found in the map.
func (u *Unit) FindPlugin(pluginID string) (plugin plugin.Plugin, found bool) {
	plugin, found = u.plugins[pluginID]
	return
}
