package core

import (
	"context"
)

// RunPlugins
func (u *Unit) RunPlugins(ctx context.Context) (err error) {
	if len(u.plugins) <= 0 {
		return
	}
	for _, plug := range u.plugins {
		if err = plug.Run(ctx); err != nil {
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
			return
		}
	}
	return
}
