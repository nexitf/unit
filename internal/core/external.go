package core

import (
	"context"
	"fmt"
	"reflect"

	"github.com/nexitf/unit/internal/core/plugin"
)

type Connector struct {
	pluginID string
	name     string
	updater  plugin.Updater
}

// LoadExternal
func (u *unit) LoadExternal(ctx context.Context) (err error) {
	if len(u.externals) <= 0 {
		return
	}
	for id, connector := range u.externals {
		plugin, exist := u.LoadPlugin(connector.pluginID)
		if !exist {
			_ = exist
		}
		err = plugin.Watch(ctx, connector.name, connector.updater)
		if err != nil {
			return err
		}
		fmt.Printf("Load: %s %s %s\n", id, plugin.Name(), connector.name)
	}
	return
}

// BindExternal
func BindExternal(varp plugin.Type, name string) {
	refVal := reflect.ValueOf(varp)
	if refVal.Kind() != reflect.Ptr || refVal.IsNil() {
		panic(ErrInvalidVariablePointer)
	}

	// Plugin ID
	pluginID := varp.PluginID()

	plugin, exist := u.LoadPlugin(pluginID)
	if !exist {
		panic(ErrVariableCanNotBeBound)
	}

	// Bind the variable
	updater := plugin.NewUpdater()
	updater.Bind(varp)

	// Bind a new external resource with a new updater
	u.externals[fmt.Sprintf("plugin://%s/%s", pluginID, name)] = Connector{pluginID: pluginID, name: name, updater: updater}
}
