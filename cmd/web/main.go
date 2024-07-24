package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"opus-classical-go/internal/config"
	"opus-classical-go/internal/models"
	"os"
)

type application struct {
	cfg       *config.Config
	periods   *models.PeriodModel
	composers *models.ComposerModel
}

func main() {
	initLogger()
	err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	cfg := config.Get()
	db, err := openDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()
	app := application{
		cfg:       cfg,
		periods:   &models.PeriodModel{DB: db},
		composers: &models.ComposerModel{DB: db},
	}

	composers, err := app.composers.GetAll()
	if err != nil {
		log.Fatalf("Failed to get composers: %v", err)
	}
	for _, period := range composers {
		fmt.Printf("%#v\n", period)
	}

	slog.Info("Web server started", "port", app.cfg.Port)
	err = http.ListenAndServe(":4000", app.routes())
	slog.Error(err.Error())
	os.Exit(1)
}
