package rpc

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type db struct {
	postgres *pgxpool.Pool
}

// New db interface
func NewDB(postgresURL string) (*db, error) {
	postgres, err := pgxpool.Connect(context.Background(), postgresURL)
	if err != nil {
		return nil, err
	}

	return &db{postgres: postgres}, nil
}
