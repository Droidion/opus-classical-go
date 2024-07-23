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
	err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	app := application{
		cfg: config.Get(),
	}
	db, err := openDB(app.cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()
	slog.Info("Web server started", "port", app.cfg.Port)
	err = http.ListenAndServe(":4000", app.routes())
	slog.Error(err.Error())
	os.Exit(1)
}
