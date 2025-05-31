// This code file is intended to be provided to third-party configuration
// implementations as an indirect way to expose configuration options.
// Third-party implementations can directly copy and use the code from this file.
package your_package_name

import (
	"github.com/nexitf/unit/plugin"
	"github.com/nexitf/unit/plugin/config"
)

// WithVariableReady binds a ready function to all configurations.
// The provided function will be called when the configuration is ready.
// It returns a BindOption that can be used in the binding process.
func WithVariableReady(fn func()) plugin.BindOption {
	return config.WithVariableReady(fn)
}

// WithVariableChange binds a change function to all configurations.
// The provided function will be called when the configuration changes.
// It returns a BindOption that can be used in the binding process.
func WithVariableChange(fn func()) plugin.BindOption {
	return config.WithVariableChange(fn)
}

// WithBeforeUpdate attaches a callback function to be invoked before the configuration update.
// The callback receives the new value and can modify it before the update takes place.
// It returns a BindOption that can be used in the binding process.
func WithBeforeUpdate(fn func(value string) string) plugin.BindOption {
	return config.WithBeforeUpdate(fn)
}

// WithAfterUpdate attaches a callback function to be invoked after the configuration update.
// The callback receives the new value and any error that occurred during the update.
// It returns a BindOption that can be used in the binding process.
func WithAfterUpdate(fn func(value string, err error)) plugin.BindOption {
	return config.WithAfterUpdate(fn)
}
