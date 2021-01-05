package db

import (
	"database/sql"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func PostgresDB(uri string) (*sql.DB, error) {
	return sql.Open("pgx", uri)
}
