package rpc

import (
	"context"

	"github.com/burstsms/mtmo-tp/backend/lib/redis"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type db struct {
	postgres *pgxpool.Pool
	redis    *redis.Connection
}

// New db interface

func NewDB(postgresURL string, redisURL string) (*db, error) {
	postgres, err := pgxpool.Connect(context.Background(), postgresURL)
	if err != nil {
		return nil, err
	}

	redis, err := redis.Connect(redisURL)
	if err != nil {
		return nil, err
	}

	return &db{postgres: postgres, redis: redis}, nil
}

type CommandTag = pgconn.CommandTag

func (db *db) Exec(sql string, args ...interface{}) (CommandTag, error) {
	return db.postgres.Exec(bg(), sql, args...)
}

func bg() context.Context {
	return context.Background()
}
