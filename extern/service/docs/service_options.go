// This code file is intended to be provided to third-party service
// implementations as an indirect way to expose service options.
// Third-party implementations can directly copy and use the code from this file.
package your_package_name

import (
	"github.com/nexitf/unit/extern/service"
	"github.com/nexitf/unit/plugin"
)

// WithVariableReady binds a ready callback function to all services.
// The provided function will be called when the service variables are ready.
func WithVariableReady(fn func()) plugin.BindOption {
	return service.WithVariableReady(fn)
}

// WithVariableChange binds a change callback function to all services.
// The provided function will be called when the service variables change.
func WithVariableChange(fn func()) plugin.BindOption {
	return service.WithVariableChange(fn)
}

// WithBeforeUpdate attaches a callback function that will be invoked before the service endpoints are updated.
// The callback receives the current list of endpoints and should return a modified list.
func WithBeforeUpdate(fn func(endpoints []service.Endpoint) []service.Endpoint) plugin.BindOption {
	return service.WithBeforeUpdate(fn)
}

// WithAfterUpdate attaches a callback function that will be invoked after the service endpoints are updated.
// The callback receives the updated list of endpoints and any error that occurred during the update process.
func WithAfterUpdate(fn func(endpoints []service.Endpoint, err error)) plugin.BindOption {
	return service.WithAfterUpdate(fn)
}

// WithTag specifies a tag to filter the service's endpoints.
// Only endpoints with the specified tag will be considered.
func WithTag(tag string) plugin.BindOption {
	return service.WithTag(tag)
}

// WithDirectAddr binds the client to the specified direct addresses.
// Each address will be treated as a service endpoint with default settings.
func WithDirectAddr(addrs ...string) plugin.BindOption {
	return service.WithDirectAddr(addrs...)
}

// WithDirectEndpoint binds the client to the specified direct endpoints.
// The provided endpoints will be used directly by the client.
func WithDirectEndpoint(endpoints ...service.Endpoint) plugin.BindOption {
	return service.WithDirectEndpoint(endpoints...)
}
