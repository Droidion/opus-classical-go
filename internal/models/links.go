package models

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Link struct {
	RecordingId int    `db:"recording_id"`
	Link        string `db:"link"`
	Streamer    string `db:"streamer"`
	Icon        string `db:"icon"`
	LinkPrefix  string `db:"link_prefix"`
}

type LinkModel struct {
	DB *pgxpool.Pool
}

func (m *LinkModel) GetLinksByRecordings(recordingIDs []int32) ([]Link, error) {
	rows, err := m.DB.Query(context.Background(), "SELECT * FROM links_with_streamers WHERE recording_id = any($1)", recordingIDs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query links")
	}
	links, err := pgx.CollectRows(rows, pgx.RowToStructByName[Link])
	if err != nil {
		return nil, errors.Wrap(err, "failed to map links")
	}
	return links, nil
}
