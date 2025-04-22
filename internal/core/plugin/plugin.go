package plugin

import (
	"context"
	"sync"
	"time"

	"github.com/nexitf/unit/internal/core/utils"
	"github.com/nexitf/unit/internal/errors"
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
	Name    string `json:"Name,omitempty"`
	Version string `json:"Version,omitempty"`
	Author  string `json:"Author,omitempty"`
	Package string `json:"Package,omitempty"`
}

type Plugin interface {

	// Name
	Name() string

	// About
	About() (about About)

	// Run
	Run(ctx context.Context) (err error)

	// Stop
	Stop(ctx context.Context) (err error)

	// Bind
	Bind(ctx context.Context, name string, updater Updater) (err error)

	// NewUpdater creates a new variable updater.
	NewUpdater() (updater Updater)
}

var (
	ErrPluginAlreadyExist = errors.New("plugin already exist")
)

var (
	mutex   sync.RWMutex
	plugins map[string]Plugin
)

func init() {
	plugins = make(map[string]Plugin)
}

// Register registers a plugin operator.
func Register(plugin Plugin) (pluginID string) {
	mutex.Lock()
	defer mutex.Unlock()
	// Generate a random string as pluginID
	pluginID, err := utils.RandomString(32)
	if err != nil {
		panic(err)
	}
	_, exist := plugins[pluginID]
	if exist {
		panic(ErrPluginAlreadyExist)
	}
	plugins[pluginID] = plugin
	return
}

// Load returns a plugin.
func Load(pluginID string) (plug Plugin, exist bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	plug, exist = plugins[pluginID]
	return
}

// LoadPlugins returns all plugins.
func LoadPlugins() map[string]Plugin {
	mutex.RLock()
	defer mutex.RUnlock()
	return plugins
}
