package core

import (
	"context"
	"fmt"
	"reflect"

	"github.com/nexitf/unit/internal/core/plugin"
)

type Connector struct {
	name     string
	pluginID string
	updater  plugin.Updater
}

// LoadExternal
func (u *unit) LoadExternal(ctx context.Context) (err error) {
	if len(u.externals) <= 0 {
		return
	}
	for _, connector := range u.externals {
		plug, exist := plugin.Load(connector.pluginID)
		if !exist {
			_ = exist
		}
		err = plug.Watch(ctx, connector.name, connector.updater)
		if err != nil {
			return err
		}
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

	plug, exist := plugin.Load(pluginID)
	if !exist {
		panic(ErrVariableCanNotBeBound)
	}

	path := fmt.Sprintf("plugin://%s/%s", pluginID, name)
	if _, ok := u.externals[path]; ok {
		panic(ErrDuplicateExternalName)
	}

	// Bind the variable
	updater := plug.NewUpdater()
	updater.Bind(varp)

	// Bind a new external resource with a new updater
	u.externals[path] = Connector{pluginID: pluginID, name: name, updater: updater}
}
