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

// String returns a string representation of the Connector.
// It constructs a URI-like string in the format "plugin://<pluginID>/<name>",
// where <pluginID> is the ID of the associated plugin and <name> is the name of the external resource.
// This string can be used for logging, debugging, or identifying the Connector instance.
//
// Returns:
//   - string: A URI-like string representing the Connector.
func (c Connector) String() string {
	return "plugin://" + c.pluginID + "/" + c.name
}

// LoadExternals loads all external resources associated with the Unit instance.
// It iterates through each external connector, finds the corresponding plugin,
// and binds the external resource to the plugin. If any binding operation fails,
// it logs the error and returns immediately. Otherwise, it logs a success message for each resource.
//
// Parameters:
//   - ctx: A context used to manage the lifecycle of the loading process,
//     which can be used for cancellation or timeout control.
//
// Returns:
//   - error: Returns nil if all external resources are loaded successfully.
//     Otherwise, returns the error encountered during the process.
func (u *Unit) LoadExternals(ctx context.Context) (err error) {
	// Associate external dependency names with the corresponding plugins
	for _, connector := range u.externals {
		spend := utils.CallSpend(func() {
			pluginID := connector.varp.PluginID()
			plugin, found := u.plugins[pluginID]
			if !found {
				err = ErrPluginNotFound
			} else {
				connector.plugin = plugin
				connector.pluginID = pluginID
				if _, ok := u.exloadeds[connector.String()]; ok {
					err = ErrDuplicateExternalName
				} else {
					u.exloadeds[connector.String()] = struct{}{}
					// Bind the variable
					connector.updater = plugin.NewUpdater()
					connector.updater.Bind(connector.varp, connector.opts...)
					err = connector.plugin.Bind(ctx, connector.name, connector.updater)
				}
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

// BindExternal binds a new external resource to the Unit instance.
// It first validates that the provided variable pointer is valid.
// If the Unit is already running, it attempts to bind the resource immediately.
// Finally, it adds the new connector to the list of external resources.
//
// Parameters:
//   - varp: The external resource to bind, must be a non - nil pointer.
//   - name: The name of the external resource.
//   - comment: A comment describing the external resource.
//   - pc: A pointer to a runpoint.PCounter.
//   - opts: Optional binding options.
//
// Returns:
//   - error: Returns nil if the binding is successful.
//     Otherwise, returns the error encountered during the process.
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
		connector.plugin = plugin
		connector.pluginID = pluginID
		if _, ok := u.exloadeds[connector.String()]; ok {
			err = ErrDuplicateExternalName
		} else {
			u.exloadeds[connector.String()] = struct{}{}
			// Bind the variable
			connector.updater = plugin.NewUpdater()
			connector.updater.Bind(varp, opts...)
			// Load now
			err = connector.plugin.Bind(context.TODO(), name, connector.updater)
			if err != nil {
				logkit.ErrorWrap(err, "bind external resource failed")
			}
		}
	}

	// Bind a new external resource
	if err == nil {
		u.externals = append(u.externals, connector)
	}

	return err
}

// BindExternal is a wrapper function that calls the BindExternal method of the Unit instance.
// It provides a convenient way to bind an external resource without accessing the Unit instance directly.
//
// Parameters:
//   - varp: The external resource to bind, must be a non - nil pointer.
//   - name: The name of the external resource.
//   - comment: A comment describing the external resource.
//   - pc: A pointer to a runpoint.PCounter.
//   - opts: Optional binding options.
//
// Returns:
//   - error: Returns nil if the binding is successful.
//     Otherwise, returns the error encountered during the process.
func BindExternal(varp plugin.Resource, name, comment string, pc *runpoint.PCounter, opts ...plugin.BindOption) (err error) {
	return u.BindExternal(varp, name, comment, pc, opts...)
}
