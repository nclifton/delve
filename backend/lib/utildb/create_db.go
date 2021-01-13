package util

import (
	"context"

	_ "github.com/lib/pq"

	"github.com/burstsms/mtmo-tp/backend/lib/logger"

	"database/sql"
	"fmt"
)

func CreateDB(host, user, password, ssl, authDB, dbName string) {
	ctx := context.Background()
	stLog := logger.NewLogger()

	connStr := fmt.Sprintf("host=%s user=%s password=%s sslmode=%s database=%s", host, user, password, ssl, authDB)
	stLog.Infof(ctx, "CreateDB", "connStr: %s", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		stLog.Fatalf(ctx, "CreateDB", "Failed to connect to pg in CreateDB: %s\n", err)
	}

	statement := fmt.Sprintf("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%s');", dbName)

	row := db.QueryRow(statement)
	var exists bool
	err = row.Scan(&exists)
	if err != nil {
		stLog.Fatalf(ctx, "CreateDB", "Failed to scan for existence of db: %s, err: %s\n", dbName, err)
	}

	if !exists {
		statement = fmt.Sprintf("CREATE DATABASE %s;", dbName)
		_, err = db.Exec(statement)
		if err != nil {
			stLog.Fatalf(ctx, "CreateDB", "Failed to CREATE DATABASE: %s, err: %s\n", dbName, err)
		}
	}
}
