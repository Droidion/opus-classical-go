package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.composersView)
	mux.HandleFunc("GET /composer/{composerSlug}", app.worksView)
	mux.HandleFunc("GET /composer/{composerSlug}/work/{workID}", app.recordingsView)
	return mux
}
