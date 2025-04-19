package core

const (
	Version = "1.0.0"
)

type RtExternal struct {
	Path     string `json:"path,omitempty"`
	Type     string `json:"type,omitempty"`
	Comment  string `json:"comment,omitempty"`
	Filename string `json:"filename,omitempty"`
	Line     int    `json:"line,omitempty"`
	Package  string `json:"package,omitempty"`
}

type Runtime struct {
	Version   string            `json:"version,omitempty"`
	Routines  []string          `json:"routines,omitempty"`
	Externals []RtExternal      `json:"externals,omitempty"`
	Plugins   map[string]string `json:"plugins,omitempty"`
}

// Inspect
func Inspect() (rt Runtime) {
	// Version
	rt.Version = Version
	// Routines
	for _, r := range u.routines {
		rt.Routines = append(rt.Routines, r.Name())
	}
	// External names
	for path, connector := range u.externals {
		rt.Externals = append(rt.Externals, RtExternal{
			Path:     path,
			Type:     connector.type_,
			Comment:  connector.comment,
			Filename: connector.pc.Filename(),
			Line:     connector.pc.Line(),
			Package:  connector.pc.PackFull(),
		})
	}
	// Plugins
	if len(u.plugins) > 0 {
		rt.Plugins = make(map[string]string)
	}
	for id, plug := range u.plugins {
		rt.Plugins[id] = plug.Name()
	}
	return
}
