package core

import (
	"context"

	"github.com/nexitf/logkit"
	"github.com/nexitf/unit/internal/core/utils"
)

// RunPlugins
func (u *Unit) RunPlugins(ctx context.Context) (err error) {
	if len(u.plugins) <= 0 {
		return
	}
	for _, plug := range u.plugins {
		spend := utils.CallSpend(func() {
			err = plug.Run(ctx)
		})
		if err != nil {
			logkit.ErrorWrap(err, "plugin run failed",
				logkit.Field("plugin", plug.Name()),
				logkit.Field("spend", spend.Seconds()),
			)
			return
		}
		logkit.Info("plugin run successfully",
			logkit.Field("plugin", plug.Name()),
			logkit.Field("spend", spend.Seconds()),
		)
	}
	return
}

// StopPlugins
func (u *Unit) StopPlugins(ctx context.Context) (err error) {
	if len(u.plugins) <= 0 {
		return
	}
	for _, plug := range u.plugins {
		spend := utils.CallSpend(func() {
			err = plug.Stop(ctx)
		})
		if err != nil {
			logkit.ErrorWrap(err, "plugin stop failed",
				logkit.Field("plugin", plug.Name()),
				logkit.Field("spend", spend.Seconds()),
			)
			return
		}
		logkit.Info("plugin stop successfully",
			logkit.Field("plugin", plug.Name()),
			logkit.Field("spend", spend.Seconds()),
		)
	}
	return
}
