package extern

import (
	"github.com/nexitf/unit/internal/core"
	"github.com/nexitf/unit/internal/core/plugin"
	"github.com/thecxx/runpoint"
)

// Bind binds a variable to a external name.
func Bind(varp plugin.Resource, name string, opts ...plugin.BindOption) (err error) {
	return core.BindExternal(varp, name, "", runpoint.PC(1), opts...)
}

// BindWithComment binds a variable to a external name.
func BindWithComment(varp plugin.Resource, name, comment string, opts ...plugin.BindOption) (err error) {
	return core.BindExternal(varp, name, comment, runpoint.PC(1), opts...)
}
