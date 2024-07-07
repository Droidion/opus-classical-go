package main

import (
	"log"
	"log/slog"
	"net/http"
	"opus-classical-go/internal/config"
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
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.composersView)
	mux.HandleFunc("GET /composer/{composerSlug}", app.worksView)
	mux.HandleFunc("GET /composer/{composerSlug}/work/{workID}", app.recordingsView)
	slog.Info("Web server started", "port", app.cfg.Port)
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
