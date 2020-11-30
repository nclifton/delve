package main

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kelseyhightower/envconfig"

	util "github.com/burstsms/mtmo-tp/backend/lib/utildb"
)

type Env struct {
	Host     string `envconfig:"POSTGRES_HOST"`
	User     string `envconfig:"POSTGRES_USER"`
	Password string `envconfig:"POSTGRES_PASSWORD"`
	SSL      string `envconfig:"POSTGRES_SSL"`
	AuthDB   string `envconfig:"POSTGRES_AUTH_DB"`
	Name     string `envconfig:"POSTGRES_NAME"`
}

func main() {
	var env Env
	err := envconfig.Process("optout", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	util.CreateDB(env.Host, env.User, env.Password, env.SSL, env.AuthDB, env.Name)

	// initialise migrations
	m, err := migrate.New(
		"file:///sql",
		fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s", env.User, env.Password, env.Host, env.Name, env.SSL),
	)
	if err != nil {
		log.Fatalf("failed to initialise golang-migrate connection: %s\n", err)
	}

	// apply migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrations: %s\n", err)
	}
}
