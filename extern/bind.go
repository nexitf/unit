package extern

import (
	"github.com/nexitf/unit/extern/plugin"
	"github.com/nexitf/unit/internal/core"
)

// Bind binds a variable to a external name.
func Bind(varp plugin.Type, name string) {
	core.BindExternal(varp, name)
}
