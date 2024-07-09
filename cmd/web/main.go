package main

import (
	"log"
	"log/slog"
	"net/http"
	"opus-classical-go/internal/config"
	"os"
)

type application struct {
	cfg *config.Config
}

func main() {
	initLogger()
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	app := application{
		cfg: config.Get(),
	}

	slog.Info("Web server started", "port", app.cfg.Port)
	err := http.ListenAndServe(":4000", app.routes())
	slog.Error(err.Error())
	os.Exit(1)
}
