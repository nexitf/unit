package plugin

import (
	"context"
	"crypto/rand"
	"errors"
	"sync"
)

// Plugin resource-binding types must implement
// the Variable interface.
type Variable interface {
	// PluginID returns the plugin ID.
	PluginID() string
}

type BindOption func(Variable) (used bool)

type Updater interface {
	// Bind binds the plugin to the variable.
	Bind(varp any, opts ...BindOption)

	// // Update updates the plugin with the given value.
	// update(...) (err error)
}

type Plugin interface {

	// Name
	Name() string

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
	pluginID, err := randomString(32)
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

const (
	charset = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// randomString generates a random string of the given length.
func randomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	for i := 0; i < length; i++ {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}
