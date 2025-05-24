package core

import (
	"context"
	"reflect"

	"github.com/nexitf/logkit"
	"github.com/nexitf/unit/internal/core/plugin"
	"github.com/nexitf/unit/internal/core/utils"
	"github.com/thecxx/runpoint"
)

type Connector struct {
	name     string
	type_    string
	varp     plugin.Resource
	opts     []plugin.BindOption
	updater  plugin.Updater // [variable pointer] <- [updater] <- [plugin]
	plugin   plugin.Plugin
	pluginID string
	comment  string
	pc       *runpoint.PCounter
}

func (c Connector) String() string {
	return "plugin://" + c.pluginID + "/" + c.name
}

// LoadExternal
func (u *Unit) LoadExternals(ctx context.Context) (err error) {
	// Associate external dependency names with the corresponding plugins
	for _, connector := range u.externals {
		spend := utils.CallSpend(func() {
			pluginID := connector.varp.PluginID()
			plugin, found := u.plugins[pluginID]
			if !found {
				err = ErrPluginNotFound
			} else {
				// Bind the variable
				connector.updater = plugin.NewUpdater()
				connector.updater.Bind(connector.varp, connector.opts...)
				connector.plugin = plugin
				connector.pluginID = pluginID
				err = connector.plugin.Bind(ctx, connector.name, connector.updater)
			}
		})
		if err != nil {
			logkit.ErrorWrap(err, "bind external resource failed",
				logkit.Field("name", connector.name),
				logkit.Field("spend", spend.Seconds()),
			)
			return err
		}
		logkit.Info("load external resource successfully",
			logkit.Field("name", connector.name),
			logkit.Field("spend", spend.Seconds()),
		)
	}
	return
}

// BindExternal
func (u *Unit) BindExternal(varp plugin.Resource, name, comment string, pc *runpoint.PCounter, opts ...plugin.BindOption) (err error) {
	refVal := reflect.ValueOf(varp)
	if refVal.Kind() != reflect.Ptr || refVal.IsNil() {
		logkit.PanicWrap(ErrInvalidVariablePointer, "varp must be a pointer")
		panic(ErrInvalidVariablePointer)
	}

	refType := reflect.TypeOf(varp).Elem()

	connector := &Connector{
		name:    name,
		varp:    varp,
		opts:    opts,
		comment: comment,
		type_:   refType.PkgPath() + "." + refType.Name(),
		pc:      pc,
	}

	if u.IsRunning() {
		// Plugin ID
		pluginID := varp.PluginID()

		plugin, found := u.plugins[pluginID]
		if !found {
			return ErrPluginNotFound
		}

		// Bind the variable
		connector.updater = plugin.NewUpdater()
		connector.updater.Bind(varp, opts...)
		connector.plugin = plugin
		connector.pluginID = pluginID

		// Load now
		err = connector.plugin.Bind(context.TODO(), name, connector.updater)
		if err != nil {
			logkit.ErrorWrap(err, "bind external resource failed")
		}
	}

	// Bind a new external resource
	u.externals = append(u.externals, connector)

	return err
}

// BindExternal
func BindExternal(varp plugin.Resource, name, comment string, pc *runpoint.PCounter, opts ...plugin.BindOption) (err error) {
	return u.BindExternal(varp, name, comment, pc, opts...)
}
