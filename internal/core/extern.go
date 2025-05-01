package core

import (
	"context"
	"reflect"

	"github.com/nexitf/logkit"
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
func (u *Unit) LoadExternals(ctx context.Context) (err error) {
	if len(u.externals) <= 0 {
		return
	}
	// Associate external dependency names with the corresponding plugins
	for _, connector := range u.externals {
		err = connector.plugin.Bind(ctx, connector.name, connector.updater)
		if err != nil {
			logkit.ErrorWrap(err, "bind external resource failed")
			return err
		}
	}
	return
}

// BindExternal
func (u *Unit) BindExternal(varp plugin.Resource, name, comment string, pc *runpoint.PCounter, opts ...plugin.BindOption) (err error) {
	refVal := reflect.ValueOf(varp)
	if refVal.Kind() != reflect.Ptr || refVal.IsNil() {
		if u.IsRunning() {
			logkit.ErrorWrap(ErrInvalidVariablePointer, "varp must be a pointer")
			return ErrInvalidVariablePointer
		} else {
			logkit.PanicWrap(ErrInvalidVariablePointer, "varp must be a pointer")
			panic(ErrInvalidVariablePointer)
		}
	}

	// Plugin ID
	pluginID := varp.PluginID()

	plug, exist := plugin.Load(pluginID)
	if !exist {
		if u.IsRunning() {
			logkit.ErrorWrap(ErrVariableCanNotBeBound, "plugin not found")
			return ErrVariableCanNotBeBound
		} else {
			logkit.PanicWrap(ErrVariableCanNotBeBound, "plugin not found")
			panic(ErrVariableCanNotBeBound)
		}
	}

	path := "plugin://" + pluginID + "/" + name
	if _, ok := u.externals[path]; ok {
		if u.IsRunning() {
			logkit.ErrorWrap(ErrDuplicateExternalName, "external resource already exists")
			return ErrDuplicateExternalName
		} else {
			logkit.PanicWrap(ErrDuplicateExternalName, "external resource already exists")
			panic(ErrDuplicateExternalName)
		}
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

	if u.IsRunning() {
		// Load now
		err = plug.Bind(context.TODO(), name, updater)
		if err != nil {
			logkit.ErrorWrap(err, "bind external resource failed")
		}
	}

	return err
}

// BindExternal
func BindExternal(varp plugin.Resource, name, comment string, pc *runpoint.PCounter, opts ...plugin.BindOption) (err error) {
	return u.BindExternal(varp, name, comment, pc, opts...)
}
