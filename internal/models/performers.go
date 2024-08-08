package models

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Performer struct {
	RecordingId int         `db:"recording_id"`
	FirstName   pgtype.Text `db:"first_name"`
	LastName    string      `db:"last_name"`
	Instrument  string      `db:"instrument"`
	Priority    pgtype.Int4 `db:"priority"`
}

type PerformerModel struct {
	DB *pgxpool.Pool
}

func (m *PerformerModel) GetPerformersByRecordings(recordingIDs []int) ([]Performer, error) {
	rows, err := m.DB.Query(context.Background(), "SELECT * FROM performers_with_instruments WHERE recording_id = any($1)", recordingIDs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query performers")
	}
	performers, err := pgx.CollectRows(rows, pgx.RowToStructByName[Performer])
	if err != nil {
		return nil, errors.Wrap(err, "failed to map performers")
	}
	return performers, nil
}
