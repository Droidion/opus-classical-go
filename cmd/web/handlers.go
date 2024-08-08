package main

import (
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

type recordingData struct {
	models.Recording
	ImagesURL  string
	Performers []models.Performer
	Links      []models.Link
}

type recordingsData struct {
	Work       *models.Work
	Composer   *models.Composer
	Recordings []recordingData
}

func (app *application) recordingsView(w http.ResponseWriter, r *http.Request) {
	workID, err := strconv.Atoi(r.PathValue("workID"))
	if err != nil || workID < 1 {
		slog.Error("Error parsing work ID", "error", err)
		http.NotFound(w, r)
		return
	}

	slug := r.PathValue("composerSlug")
	if slug == "" {
		slog.Error("Error parsing composer slug", "error", err)
		http.NotFound(w, r)
		return
	}

	work, err := app.works.GetWorkByID(workID)
	if err != nil {
		slog.Error("Error finding work with given ID", "error", err)
		http.NotFound(w, r)
		return
	}

	composer, err := app.composers.GetOneBySlug(&slug)
	if err != nil {
		slog.Error("Error finding composer by slug", "error", err)
		http.NotFound(w, r)
		return
	}

	recordings, err := app.recordings.GetRecordingsByWork(workID)
	if err != nil {
		slog.Error("Error finding recordings by work ID", "error", err)
		http.NotFound(w, r)
		return
	}

	recordingsIDs := lo.Map(recordings, func(record models.Recording, _ int) int {
		return record.ID
	})

	performers, err := app.performers.GetPerformersByRecordings(recordingsIDs)
	if err != nil {
		slog.Error("Error finding performers for recordings", "error", err)
		http.NotFound(w, r)
		return
	}

	links, err := app.links.GetLinksByRecordings(recordingsIDs)
	if err != nil {
		slog.Error("Error finding links for recordings", "error", err)
		http.NotFound(w, r)
		return
	}

	recordingData := lo.Map(recordings, func(record models.Recording, _ int) recordingData {
		return recordingData{
			Recording: record,
			ImagesURL: app.cfg.ImagesURL,
			Performers: lo.Filter(performers, func(performer models.Performer, _ int) bool {
				return performer.RecordingId == record.ID
			}),
			Links: lo.Filter(links, func(link models.Link, _ int) bool {
				return link.RecordingId == record.ID
			}),
		}
	})

	recordData := recordingsData{
		Work:       work,
		Composer:   composer,
		Recordings: recordingData,
	}

	app.render(w, r, http.StatusOK, "recordings.gohtml", recordData)
}
