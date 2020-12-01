package util

import (
	_ "github.com/lib/pq"

	"database/sql"
	"fmt"
	"log"
)

func CreateDB(host, user, password, ssl, authDB, dbName string) {
	connStr := fmt.Sprintf("host=%s user=%s password=%s sslmode=%s database=%s", host, user, password, ssl, authDB)
	log.Println("CreateDB connStr", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to pg in CreateDB: %s\n", err)
	}

	statement := fmt.Sprintf("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%s');", dbName)

	row := db.QueryRow(statement)
	var exists bool
	err = row.Scan(&exists)
	if err != nil {
		log.Fatalf("Failed to scan for existence of db: %s, err: %s\n", dbName, err)
	}

	if !exists {
		statement = fmt.Sprintf("CREATE DATABASE %s;", dbName)
		_, err = db.Exec(statement)
		if err != nil {
			log.Fatalf("Failed to CREATE DATABASE: %s, err: %s\n", dbName, err)
		}
	}
}
