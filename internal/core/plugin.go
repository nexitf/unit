package core

import (
	"context"
)

// RunPlugin
func (u *unit) RunPlugin(ctx context.Context) (err error) {
	for _, plug := range u.plugins {
		if err = plug.Run(ctx); err != nil {
			return
		}
	}
	return
}

// StopPlugin
func (u *unit) StopPlugin(ctx context.Context) (err error) {
	for _, plug := range u.plugins {
		if err = plug.Stop(ctx); err != nil {
			return
		}
	}
	return
}
