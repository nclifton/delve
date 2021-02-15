package main

import (
	util "github.com/burstsms/mtmo-tp/backend/lib/utildb"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	util.Migrate()
}
