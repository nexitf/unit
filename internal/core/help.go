package core

import (
	"fmt"
	"sort"
	"time"
)

type RtRoutine struct {
	Name     string `json:"Name,omitempty"`
	RunTime  int64  `json:"RunTime,omitempty"`
	StopTime int64  `json:"StopTime,omitempty"`
	Time     string `json:"Time,omitempty"`
	Error    string `json:"Error,omitempty"`
}

type RtSnapshot struct {
	Time string `json:"Time"`
	Data string `json:"Data"`
}

type RtExternal struct {
	Path     string     `json:"Path,omitempty"`
	Var      string     `json:"Var,omitempty"`
	Snapshot RtSnapshot `json:"Snapshot,omitempty"`
	Type     string     `json:"Type,omitempty"`
	Comment  string     `json:"Comment,omitempty"`
	Filename string     `json:"Filename,omitempty"`
	Line     int        `json:"Line,omitempty"`
	Package  string     `json:"Package,omitempty"`
}

type RtPlugin struct {
	PluginID    string `json:"PluginID,omitempty"`
	Name        string `json:"Name,omitempty"`
	Version     string `json:"Version,omitempty"`
	Description string `json:"Description,omitempty"`
	Package     string `json:"Package,omitempty"`
}

type Runtime struct {
	Version   string       `json:"Version"`
	Fatal     string       `json:"Fatal,omitempty"`
	Running   bool         `json:"Running"`
	Routines  []RtRoutine  `json:"Routines"`
	Externals []RtExternal `json:"Externals"`
	Plugins   []RtPlugin   `json:"Plugins"`
}

// Inspect
func Inspect() (rt Runtime) {
	// Version
	rt.Version = Version
	// Error
	if u.fatal != nil {
		rt.Fatal = u.fatal.Error()
	}
	// Status
	rt.Running = u.IsRunning()
	// Routines
	rt.Routines = make([]RtRoutine, 0)
	for _, l := range u.scheduler.Routines() {
		v := RtRoutine{
			Name: l.Name(),
		}
		if rt := l.RunTime(); !rt.IsZero() {
			v.RunTime = rt.UnixMilli()
		}
		if st := l.StopTime(); !st.IsZero() {
			v.StopTime = st.UnixMilli()
		}
		if v.RunTime > 0 {
			if v.StopTime <= 0 {
				v.Time = formatDurationMS(time.Now().UnixMilli() - v.RunTime)
			} else {
				v.Time = formatDurationMS(v.StopTime - v.RunTime)
			}
		}
		if err := l.Err(); err != nil {
			v.Error = err.Error()
		}
		rt.Routines = append(rt.Routines, v)
	}
	// External names
	rt.Externals = make([]RtExternal, 0)
	for _, connector := range u.externals {
		v := RtExternal{
			Path:     connector.String(),
			Var:      fmt.Sprintf("%p", connector.varp),
			Type:     connector.type_,
			Comment:  connector.comment,
			Filename: connector.pc.Filename(),
			Line:     connector.pc.Line(),
			Package:  connector.pc.PackFull(),
		}
		if connector.updater != nil {
			snap := connector.updater.Snapshot()
			if !snap.Time.IsZero() {
				v.Snapshot.Time = snap.Time.In(time.Local).Format("2006-01-02 15:04:05.000000")
				v.Snapshot.Data = snap.Data
			}
		}
		rt.Externals = append(rt.Externals, v)
	}
	// Plugins
	rt.Plugins = make([]RtPlugin, 0)
	for pluginID, plugin := range u.plugins {
		about := plugin.About()
		v := RtPlugin{
			PluginID:    pluginID,
			Name:        about.Name,
			Version:     about.Version,
			Description: about.Description,
			Package:     about.Package,
		}
		rt.Plugins = append(rt.Plugins, v)
	}
	// Sort plugins
	sort.Slice(rt.Plugins, func(i, j int) bool {
		return rt.Plugins[i].Name < rt.Plugins[j].Name
	})
	return
}

// formatDurationMS converts a duration in milliseconds (int64)
// into a human-readable string format: "Xd Xh Xm Xs"
func formatDurationMS(ms int64) string {
	if ms < 0 {
		// Handle negative durations
		return "-" + formatDurationMS(-ms)
	}

	// Convert milliseconds to time.Duration
	duration := time.Duration(ms) * time.Millisecond

	days := duration / (24 * time.Hour)
	duration %= 24 * time.Hour

	hours := duration / time.Hour
	duration %= time.Hour

	minutes := duration / time.Minute
	duration %= time.Minute

	seconds := duration / time.Second
	duration %= time.Second

	// Build formatted string
	result := ""
	if days > 0 {
		result += fmt.Sprintf("%dd", days)
	}
	if hours > 0 || len(result) > 0 {
		result += fmt.Sprintf("%02dh", hours)
	}
	if minutes > 0 || len(result) > 0 {
		result += fmt.Sprintf("%02dm", minutes)
	}
	if seconds > 0 || len(result) > 0 {
		result += fmt.Sprintf("%02ds", seconds)
	}

	return result
}
