package util

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/kelseyhightower/envconfig"

	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/nr"
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

func Migrate(dbName string) {
	ctx := context.Background()
	stLog := logger.NewLogger()

	stLog.Info(ctx, "Migrate", "Starting service...")

	var env Env
	err := envconfig.Process(dbName, &env)
	if err != nil {
		stLog.Fatalf(ctx, "Migrate", "Failed to read env vars: %s", err)
	}

	stLog.Infof(ctx, "Migrate", "ENV: %+v", env)

	// Register service with New Relic
	nr.CreateApp(&nr.Options{
		AppName:                  env.NRName,
		NewRelicLicense:          env.NRLicense,
		DistributedTracerEnabled: env.NRTracing,
	})

	stLog.Info(ctx, "Migrate", "Service started")

	CreateDB(env.Host, env.User, env.Password, env.SSL, env.AuthDB, env.Name)

	// initialise migrations
	m, err := migrate.New(
		"file:///sql",
		fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s", env.User, env.Password, env.Host, env.Name, env.SSL),
	)
	if err != nil {
		stLog.Fatalf(ctx, "Migrate", "Failed to initialise golang-migrate connection: %s\n", err)
	}

	// apply migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		stLog.Fatalf(ctx, "Migrate", "Failed to run migrations: %s\n", err)
	}

	stLog.Info(ctx, "Migrate", "Migration completed successfully!")
}
