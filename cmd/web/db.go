package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func openDB(databaseURL string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}
	err = dbpool.Ping(context.Background())
	if err != nil {
		dbpool.Close()
		return nil, err
	}
	return dbpool, nil
}
