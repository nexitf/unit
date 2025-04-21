package core

import (
	"fmt"
	"time"

	"github.com/nexitf/unit/internal/core/plugin"
)

const (
	Version = "1.0.0"
)

type RtRoutine struct {
	Name     string `json:"Name"`
	RunTime  int64  `json:"RunTime"`
	StopTime int64  `json:"StopTime"`
	Error    string `json:"Error"`
}

type RtSnapshot struct {
	Time string `json:"Time"`
	Data string `json:"Data"`
}

type RtExternal struct {
	Path     string     `json:"Path"`
	Var      string     `json:"Var"`
	Snapshot RtSnapshot `json:"Snapshot"`
	Type     string     `json:"Type"`
	Comment  string     `json:"Comment"`
	Filename string     `json:"Filename"`
	Line     int        `json:"Line"`
	Package  string     `json:"Package"`
}

type Runtime struct {
	Version   string                  `json:"Version"`
	Fatal     string                  `json:"Fatal"`
	Routines  []RtRoutine             `json:"Routines"`
	Externals []RtExternal            `json:"Externals"`
	Plugins   map[string]plugin.About `json:"Plugins"`
}

// Inspect
func Inspect() (rt Runtime) {
	// Version
	rt.Version = Version
	// Error
	if err := u.Fatal(); err != nil {
		rt.Fatal = err.Error()
	}
	// Routines
	for _, r := range u.routines {
		l := r.(*Launcher)
		v := RtRoutine{
			Name: r.Name(),
		}
		if !l.runTime.IsZero() {
			v.RunTime = l.runTime.UnixMilli()
		}
		if !l.stopTime.IsZero() {
			v.StopTime = l.stopTime.UnixMilli()
		}
		if l.err != nil {
			v.Error = l.err.Error()
		}
		rt.Routines = append(rt.Routines, v)
	}
	// External names
	for path, connector := range u.externals {
		v := RtExternal{
			Path:     path,
			Var:      fmt.Sprintf("%p", connector.varp),
			Type:     connector.type_,
			Comment:  connector.comment,
			Filename: connector.pc.Filename(),
			Line:     connector.pc.Line(),
			Package:  connector.pc.PackFull(),
		}
		snap := connector.updater.Snapshot()
		if !snap.Time.IsZero() {
			v.Snapshot.Time = snap.Time.In(time.Local).Format("2006-01-02 15:04:05.000000")
			v.Snapshot.Data = snap.Data
		}
		rt.Externals = append(rt.Externals, v)
	}
	// Plugins
	if len(u.plugins) > 0 {
		rt.Plugins = make(map[string]plugin.About)
	}
	for id, plug := range u.plugins {
		rt.Plugins[id] = plug.About()
	}
	return
}
