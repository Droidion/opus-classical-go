package models

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"sort"
)

type Work struct {
	ID               int         `db:"id"`
	Title            string      `db:"title"`
	YearStart        pgtype.Int4 `db:"year_start"`
	YearFinish       pgtype.Int4 `db:"year_finish"`
	AverageMinutes   pgtype.Int4 `db:"average_minutes"`
	CatalogueName    pgtype.Text `db:"catalogue_name"`
	CatalogueNumber  pgtype.Int4 `db:"catalogue_number"`
	CataloguePostfix pgtype.Text `db:"catalogue_postfix"`
	No               pgtype.Int4 `db:"no"`
	Nickname         pgtype.Text `db:"nickname"`
	ComposerID       int         `db:"composer_id"`
	Sort             pgtype.Int4 `db:"sort"`
	GenreID          int         `db:"genre_id"`
	GenreName        string      `db:"genre_name"`
}

type WorkByGenre struct {
	GenreID   int
	GenreName string
	Works     []Work
}

type WorkModel struct {
	DB *pgxpool.Pool
}

func (m *WorkModel) GetWorkByID(workID int) (*Work, error) {
	rows, err := m.DB.Query(context.Background(), "SELECT * FROM works_with_genres WHERE id = $1", workID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query work by id")
	}
	work, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Work])
	if err != nil {
		return nil, errors.Wrap(err, "failed to map work")
	}
	return &work, nil
}

func (m *WorkModel) GetWorksByComposerID(composerID int) ([]WorkByGenre, error) {
	rows, err := m.DB.Query(context.Background(), "SELECT * FROM works_with_genres WHERE composer_id = $1 ORDER BY genre_name, sort, year_finish ", composerID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query composer by slug")
	}
	works, err := pgx.CollectRows(rows, pgx.RowToStructByName[Work])
	if err != nil {
		return nil, errors.Wrap(err, "failed to map composer")
	}
	groupedWorks := lo.GroupBy(works, func(w Work) int {
		return w.GenreID
	})
	worksByGenre := lo.MapToSlice(groupedWorks, func(genreID int, works []Work) WorkByGenre {
		return WorkByGenre{
			GenreID:   genreID,
			GenreName: works[0].GenreName,
			Works:     works,
		}
	})
	sort.Slice(worksByGenre, func(i, j int) bool {
		return worksByGenre[i].GenreName < worksByGenre[j].GenreName
	})
	return worksByGenre, nil
}
