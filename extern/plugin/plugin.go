package plugin

import (
	"github.com/nexitf/unit/internal/core/plugin"
)

type Plugin = plugin.Plugin
type Variable = plugin.Variable
type BindOption = plugin.BindOption
type Updater = plugin.Updater

// Register register a plugin.
func Register(plug Plugin) (pluginID string) {
	return plugin.Register(plug)
}
