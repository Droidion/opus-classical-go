package models

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Recording struct {
	ID         int         `db:"id"`
	CoverName  string      `db:"cover_name"`
	Length     pgtype.Int4 `db:"length"`
	Label      string      `db:"label"`
	WorkId     int         `db:"work_id"`
	YearStart  pgtype.Int4 `db:"year_start"`
	YearFinish pgtype.Int4 `db:"year_finish"`
}

type RecordingModel struct {
	DB *pgxpool.Pool
}

func (m *RecordingModel) GetRecordingsByWork(workID int) ([]Recording, error) {
	rows, err := m.DB.Query(context.Background(), "SELECT * FROM recordings_with_labels WHERE work_id = $1", workID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query recordings")
	}
	recordings, err := pgx.CollectRows(rows, pgx.RowToStructByName[Recording])
	if err != nil {
		return nil, errors.Wrap(err, "failed to map recordings")
	}
	return recordings, nil
}
