package unit

import (
	"encoding/json"

	"github.com/nexitf/unit/internal/core"
)

type Runtime = core.Runtime

// Inspect returns all infomation of the current unit.
func Inspect() (rt Runtime) {
	return core.Inspect()
}

// InspectJSON returns all infomation of the current unit in JSON format.
func InspectJSON() string {
	rt := core.Inspect()
	out, err := json.MarshalIndent(rt, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(out)
}
