package plugin

import (
	"github.com/nexitf/unit/internal/core/plugin"
)

type Plugin = plugin.Plugin
type Updater = plugin.Updater
type Type = plugin.Type

// Register register a plugin.
func Register(plug Plugin) (pluginID string) {
	return plugin.Register(plug)
}
