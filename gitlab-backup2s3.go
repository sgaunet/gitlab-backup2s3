package main

import (
	"log/slog"
	"os"

	"github.com/sgaunet/gitlab-backup2s3/pkg/app"
	"github.com/sgaunet/gitlab-backup2s3/pkg/logger"
)

func main() {
	// initialize logger
	debugLevel := os.Getenv("DEBUGLEVEL")
	logger := logger.NewLogger(debugLevel)

	// initialize app
	app := app.NewApp()
	app.SetLogger(logger)

	// run app
	err := app.Run()
	if err != nil {
		logger.Error("Error executing gitlab-backup", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
