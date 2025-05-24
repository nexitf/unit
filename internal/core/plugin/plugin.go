package plugin

import (
	"context"
	"time"
)

// Plugin resource-binding types must implement
// the Resource interface.
type Resource interface {
	// PluginID returns the plugin ID.
	PluginID() string
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

	// Name
	Name() string

	// About
	About() (about About)

	// Run
	Run(ctx context.Context, pluginID string) (err error)

	// Stop
	Stop(ctx context.Context) (err error)

	// Bind
	Bind(ctx context.Context, name string, updater Updater) (err error)

	// NewUpdater creates a new variable updater.
	NewUpdater() (updater Updater)
}
