package rpc

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	AccountTable = "account"
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

type CommandTag = pgconn.CommandTag

func (db *db) Exec(sql string, args ...interface{}) (CommandTag, error) {
	return db.postgres.Exec(bg(), sql, args...)
}

func (db *db) FindByAPIKey(key string) (*Account, error) {

	var account Account

	sql := `
SELECT account.id, account.name, account.created_at, account.updated_at 
FROM account
LEFT JOIN account_api_keys as api_keys ON account.id = api_keys.account_id
WHERE api_keys.key = $1
	`

	row := db.postgres.QueryRow(bg(), sql, key)
	err := row.Scan(&account.ID, &account.Name, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func bg() context.Context {
	return context.Background()
}
