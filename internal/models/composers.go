package models

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Composer struct {
	ID            int         `db:"id"`
	FirstName     string      `db:"first_name"`
	LastName      string      `db:"last_name"`
	YearBorn      int         `db:"year_born"`
	YearDied      pgtype.Int4 `db:"year_died"`
	PeriodID      int         `db:"period_id"`
	Slug          string      `db:"slug"`
	WikipediaLink pgtype.Text `db:"wikipedia_link"`
	ImslpLink     pgtype.Text `db:"imslp_link"`
	Enabled       bool        `db:"enabled"`
	Countries     string      `db:"countries"`
}

type ComposerModel struct {
	DB *pgxpool.Pool
}

func (m *ComposerModel) GetAll() ([]Composer, error) {
	rows, err := m.DB.Query(context.Background(), "SELECT * FROM composers_with_countries")
	if err != nil {
		return nil, errors.Wrap(err, "failed to query composers")
	}
	composers, err := pgx.CollectRows(rows, pgx.RowToStructByName[Composer])
	if err != nil {
		return nil, errors.Wrap(err, "failed to map composers")
	}
	return composers, nil
}
