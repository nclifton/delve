package main

import (
	util "github.com/burstsms/mtmo-tp/backend/lib/utildb"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

// TODO: refactor this whole thing into a reusable docker service
func main() {
	util.Migrate()
}
