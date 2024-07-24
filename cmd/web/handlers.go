package main

import (
	"fmt"
	"github.com/samber/lo"
	"html/template"
	"log/slog"
	"net/http"
	"opus-classical-go/internal/models"
	"strconv"
)

type composerData struct {
	models.Period
	Composers []models.Composer
}

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

	periods, err := app.periods.GetAll()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	composers, err := app.composers.GetAll()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := lo.Map(periods, func(period models.Period, _ int) composerData {
		return composerData{
			Period: period,
			Composers: lo.Filter(composers, func(composer models.Composer, _ int) bool {
				return composer.PeriodID == period.ID
			}),
		}
	})

	err = ts.ExecuteTemplate(w, "base", data)
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
