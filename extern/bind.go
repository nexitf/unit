package extern

import (
	"github.com/nexitf/unit/extern/plugin"
	"github.com/nexitf/unit/internal/core"
	"github.com/thecxx/runpoint"
)

// Bind binds a variable to a external name.
func Bind(varp plugin.Variable, name string, opts ...plugin.BindOption) {
	core.BindExternal(varp, name, "", runpoint.PC(1), opts...)
}

// BindWithComment binds a variable to a external name.
func BindWithComment(varp plugin.Variable, name, comment string, opts ...plugin.BindOption) {
	core.BindExternal(varp, name, comment, runpoint.PC(1), opts...)
}
