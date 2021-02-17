package rpc

import (
	"context"

	"github.com/burstsms/mtmo-tp/backend/account/rpc/types"
	"github.com/burstsms/mtmo-tp/backend/lib/redis"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	AccountTable = "account"
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

	err = redis.EnableCache()
	if err != nil {
		return nil, err
	}

	return &db{postgres: postgres, redis: redis}, nil
}

type CommandTag = pgconn.CommandTag

func (db *db) Exec(sql string, args ...interface{}) (CommandTag, error) {
	return db.postgres.Exec(bg(), sql, args...)
}

func (db *db) FindByAPIKey(key string) (*types.Account, error) {

	var account types.Account

	sql := `
		SELECT a.id, a.name, a.created_at, a.updated_at, a.alaris_username, a.alaris_password, a.alaris_url
		FROM account a
		LEFT JOIN account_api_keys as ak ON a.id = ak.account_id
		WHERE ak.key = $1;
	`

	row := db.postgres.QueryRow(bg(), sql, key)
	if err := row.Scan(
		&account.ID,
		&account.Name,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.AlarisUsername,
		&account.AlarisPassword,
		&account.AlarisURL,
	); err != nil {
		return nil, err
	}

	return &account, nil
}

func (db *db) FindByID(id string) (*types.Account, error) {

	var account types.Account

	sql := `
		SELECT id, name, created_at, updated_at, alaris_username, alaris_password, alaris_url
		FROM account
		WHERE id = $1;
	`

	row := db.postgres.QueryRow(bg(), sql, id)
	if err := row.Scan(
		&account.ID,
		&account.Name,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.AlarisUsername,
		&account.AlarisPassword,
		&account.AlarisURL,
	); err != nil {
		return nil, err
	}

	return &account, nil
}

func bg() context.Context {
	return context.Background()
}
