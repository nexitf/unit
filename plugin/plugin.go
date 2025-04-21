package plugin

import (
	"github.com/nexitf/unit/internal/core/plugin"
)

type Plugin = plugin.Plugin
type About = plugin.About
type Resource = plugin.Resource
type Snapshot = plugin.Snapshot
type BindOption = plugin.BindOption
type Updater = plugin.Updater

// Register register a plugin.
func Register(plug Plugin) (pluginID string) {
	return plugin.Register(plug)
}
