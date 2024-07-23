package models

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Period struct {
	ID        int         `db:"id"`
	Name      string      `db:"name"`
	YearStart int         `db:"year_start"`
	YearEnd   pgtype.Int4 `db:"year_end"`
	Slug      string      `db:"slug"`
}

type PeriodModel struct {
	DB *pgxpool.Pool
}

func (m *PeriodModel) GetAll() ([]Period, error) {
	rows, err := m.DB.Query(context.Background(), "SELECT * FROM periods")
	if err != nil {
		return nil, err
	}
	periods, err := pgx.CollectRows(rows, pgx.RowToStructByName[Period])
	if err != nil {
		return nil, err
	}
	return periods, nil
}
