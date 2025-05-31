package plugin

import (
	"context"
	"time"
)

// Plugin resource-binding types must implement
// the Resource interface.
type Resource interface {
	// PluginID returns the plugin ID and name.
	PluginID() (pluginID string, pluginName string)
}

type Snapshot struct {
	Time time.Time
	Data string
}

type BindOption func(Resource) (used bool)

type Updater interface {
	// Bind binds the plugin to the variable.
	Bind(varp Resource, opts ...BindOption)

	// Snapshot returns the snapshot of the updater.
	Snapshot() (snapshot Snapshot)
}

type About struct {
	Name        string `json:"Name,omitempty"`
	Version     string `json:"Version,omitempty"`
	Description string `json:"Description,omitempty"`
	Package     string `json:"Package,omitempty"`
}

type Plugin interface {

	// Name returns the name of the plugin.
	Name() string

	// About provides descriptive information about the plugin.
	About() (about About)

	// Run starts the plugin's execution logic.
	// The pluginID can be used to identify or configure the plugin instance.
	Run(ctx context.Context, pluginID string) (err error)

	// Stop gracefully stops the plugin's execution.
	Stop(ctx context.Context) (err error)

	// Bind links a variable or resource to the plugin.
	// The 'name' specifies what is being bound, and 'updater' handles updates.
	Bind(ctx context.Context, name string, updater Updater) (err error)

	// NewUpdater creates and returns a new variable updater instance.
	NewUpdater() (updater Updater)
}
