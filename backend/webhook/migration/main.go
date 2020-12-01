package main

import (
	"fmt"
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	util "github.com/burstsms/mtmo-tp/backend/lib/utildb"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	Host     string `envconfig:"POSTGRES_HOST"`
	User     string `envconfig:"POSTGRES_USER"`
	Password string `envconfig:"POSTGRES_PASSWORD"`
	SSL      string `envconfig:"POSTGRES_SSL"`
	AuthDB   string `envconfig:"POSTGRES_AUTH_DB"`
	Name     string `envconfig:"POSTGRES_NAME"`

	NRName    string `envconfig:"NR_NAME"`
	NRLicense string `envconfig:"NR_LICENSE"`
	NRTracing bool   `envconfig:"NR_TRACING"`
}

func main() {
	log.Println("Starting service...")

	var env Env
	err := envconfig.Process("webhook", &env)
	if err != nil {
		log.Fatal("Failed to read env vars:", err)
	}

	log.Printf("ENV: %+v", env)

	// Register service with New Relic
	nr.CreateApp(&nr.Options{
		AppName:                  env.NRName,
		NewRelicLicense:          env.NRLicense,
		DistributedTracerEnabled: env.NRTracing,
	})

	log.Println("Service started")

	util.CreateDB(env.Host, env.User, env.Password, env.SSL, env.AuthDB, env.Name)

	// initialise migrations
	m, err := migrate.New(
		"file:///sql",
		fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s", env.User, env.Password, env.Host, env.Name, env.SSL),
	)
	if err != nil {
		log.Fatalf("Failed to initialise golang-migrate connection: %s\n", err)
	}

	// apply migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %s\n", err)
	}

	log.Println("Migration completed successfully!")
}
