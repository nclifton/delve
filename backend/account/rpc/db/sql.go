package db

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type sqlDB struct {
	sql *pgxpool.Pool
}

func NewSQLDB(db *pgxpool.Pool) DB {
	return &sqlDB{
		sql: db,
	}
}
