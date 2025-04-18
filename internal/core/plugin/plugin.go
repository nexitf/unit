package plugin

import (
	"context"
	"crypto/rand"
	"errors"
	"sync"
)

type Updater interface {
	// Bind binds the plugin to the variable.
	Bind(varp any)

	// // Update updates the plugin with the given value.
	// Update(value string) (err error)
}

type Type interface {
	// PluginID returns the plugin id.
	PluginID() string
}

type Plugin interface {

	// Name
	Name() string

	// // Run
	// Run(ctx context.Context) (err error)

	// // Stop
	// Stop(ctx context.Context) (err error)

	// Watch
	Watch(ctx context.Context, name string, updater Updater) (err error)

	// NewUpdater creates a new variable updater.
	NewUpdater() (updater Updater)
}

var (
	ErrOperatorAlreadyExist = errors.New("operator already exist")
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
		panic(ErrOperatorAlreadyExist)
	}
	plugins[pluginID] = plugin
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
