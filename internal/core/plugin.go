package core

import (
	"context"

	"github.com/nexitf/logkit"
)

// RunPlugins
func (u *Unit) RunPlugins(ctx context.Context) (err error) {
	if len(u.plugins) <= 0 {
		return
	}
	for _, plug := range u.plugins {
		if err = plug.Run(ctx); err != nil {
			logkit.ErrorWrap(err, "plugin run failed", logkit.Field("plugin", plug.Name()))
			return
		}
	}
	return
}

// StopPlugins
func (u *Unit) StopPlugins(ctx context.Context) (err error) {
	if len(u.plugins) <= 0 {
		return
	}
	for _, plug := range u.plugins {
		if err = plug.Stop(ctx); err != nil {
			logkit.ErrorWrap(err, "plugin stop failed", logkit.Field("plugin", plug.Name()))
			return
		}
	}
	return
}
