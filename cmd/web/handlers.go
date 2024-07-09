package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
)

func (app *application) composersView(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/partials/header.gohtml",
		"./ui/html/partials/footer.gohtml",
		"./ui/html/base.gohtml",
		"./ui/html/pages/composers.gohtml",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) worksView(w http.ResponseWriter, r *http.Request) {
	composerSlug := r.PathValue("composerSlug")
	if composerSlug == "" {
		slog.Error("Composer slug is empty")
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Request composer %s", composerSlug)
}

func (app *application) recordingsView(w http.ResponseWriter, r *http.Request) {
	composerSlug := r.PathValue("composerSlug")
	workID, err := strconv.Atoi(r.PathValue("workID"))
	if err != nil || workID < 1 {
		slog.Error("Error parsing work ID", "error", err)
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Request composer %s and work ID %d", composerSlug, workID)
}
