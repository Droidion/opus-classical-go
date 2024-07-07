package main

import (
	"log/slog"
	"os"
)

func initLogger() {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})

	logger := slog.New(logHandler)
	slog.SetDefault(logger)
}
