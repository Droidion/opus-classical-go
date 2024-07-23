package main

import (
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
	"time"
)

func getJSONHandler() slog.Handler {
	return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})
}

func getTintedHandler() slog.Handler {
	return tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
	})
}

func initLogger() {
	logger := slog.New(getTintedHandler())
	slog.SetDefault(logger)
}
