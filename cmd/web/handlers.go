package main

import (
	"fmt"
	"github.com/samber/lo"
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

	app.render(w, r, http.StatusOK, "composers.gohtml", data)
}

type workData struct {
	Works    []models.WorkByGenre
	Composer models.Composer
}

func (app *application) worksView(w http.ResponseWriter, r *http.Request) {
	composerSlug := r.PathValue("composerSlug")
	if composerSlug == "" {
		slog.Error("Composer slug is empty")
		http.NotFound(w, r)
		return
	}

	composer, err := app.composers.GetOneBySlug(&composerSlug)
	if err != nil {
		slog.Error("Error finding composer by slug", "error", err)
		http.NotFound(w, r)
		return
	}

	works, err := app.works.GetWorksByComposerID(composer.ID)
	if err != nil {
		slog.Error("Error finding works by composer ID", "error", err)
		http.NotFound(w, r)
		return
	}

	app.render(w, r, http.StatusOK, "works.gohtml", workData{
		Works:    works,
		Composer: *composer,
	})
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
