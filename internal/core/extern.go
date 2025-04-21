package core

import (
	"context"
	"reflect"

	"github.com/nexitf/unit/internal/core/plugin"
	"github.com/thecxx/runpoint"
)

type Connector struct {
	name    string
	type_   string
	varp    plugin.Resource
	updater plugin.Updater // [variable pointer] <- [updater] <- [plugin]
	plugin  plugin.Plugin
	comment string
	pc      *runpoint.PCounter
}

// LoadExternal
func (u *Unit) LoadExternal(ctx context.Context) (err error) {
	if len(u.externals) <= 0 {
		return
	}
	// Associate external dependency names with the corresponding plugins
	for _, connector := range u.externals {
		err = connector.plugin.Bind(ctx, connector.name, connector.updater)
		if err != nil {
			return err
		}
	}
	return
}

// BindExternal
func (u *Unit) BindExternal(varp plugin.Resource, name, comment string, pc *runpoint.PCounter, opts ...plugin.BindOption) {
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

	path := "plugin://" + pluginID + "/" + name
	if _, ok := u.externals[path]; ok {
		panic(ErrDuplicateExternalName)
	}

	// Bind the variable
	updater := plug.NewUpdater()
	updater.Bind(varp, opts...)

	refType := reflect.TypeOf(varp).Elem()

	// Bind a new external resource with a new updater
	u.externals[path] = Connector{
		name:    name,
		varp:    varp,
		updater: updater,
		plugin:  plug,
		comment: comment,
		type_:   refType.PkgPath() + "." + refType.Name(),
		pc:      pc,
	}
}

// BindExternal
func BindExternal(varp plugin.Resource, name, comment string, pc *runpoint.PCounter, opts ...plugin.BindOption) {
	u.BindExternal(varp, name, comment, pc, opts...)
}
