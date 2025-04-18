package core

import (
	"context"
)

// RunPlugins
func (u *unit) RunPlugins(ctx context.Context) (err error) {
	for _, plug := range u.plugins {
		if err = plug.Run(ctx); err != nil {
			return
		}
	}
	return
}

// StopPlugins
func (u *unit) StopPlugins(ctx context.Context) (err error) {
	for _, plug := range u.plugins {
		if err = plug.Stop(ctx); err != nil {
			return
		}
	}
	return
}
