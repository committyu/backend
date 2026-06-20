package main

import (
	"context"

	"backend/internal/pkg/logger"
	"backend/internal/usecase/game"

	"github.com/robfig/cron/v3"
)

func StartCron(resetPlayTimeUc *game.ResetPlayTime) {
	c := cron.New()

	_, err := c.AddFunc("* * * * *", func() {
		logger.Info("reset play_time cron start")

		if err := resetPlayTimeUc.Execute(context.Background()); err != nil {
			logger.Error("reset play_time cron failed", "error", err)
			return
		}

		logger.Info("reset play_time cron end")
	})

	if err != nil {
		logger.Error("cron setup failed", "error", err)
		return
	}

	c.Start()
}