package db_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	util "github.com/burstsms/mtmo-tp/backend/lib/utildb"
	"github.com/burstsms/mtmo-tp/backend/webhook/db"
	"github.com/kelseyhightower/envconfig"
)

const (
	migrationsURL = "../migration/sql"
)

var (
	testUtil  *util.TestUtil
	webhookDB *db.DB
)

type Env struct {
	Host     string `envconfig:"POSTGRES_HOST"`
	User     string `envconfig:"POSTGRES_USER"`
	Password string `envconfig:"POSTGRES_PASSWORD"`
	SSL      string `envconfig:"POSTGRES_SSL"`
	AuthDB   string `envconfig:"POSTGRES_AUTH_DB"`
	TestName string `envconfig:"POSTGRES_TEST_NAME"`
}

func bg() context.Context {
	return context.Background()
}

func TestMain(main *testing.M) {
	var env Env
	err := envconfig.Process("webhook", &env)
	if err != nil {
		log.Fatalf("failed to read env vars: %s\n", err)
	}

	util.CreateDB(env.Host, env.User, env.Password, env.SSL, env.AuthDB, env.TestName)

	pgUrl := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s", env.User, env.Password, env.Host, env.TestName, env.SSL)

	testUtil = util.NewTestUtil(pgUrl, migrationsURL)

	webhookDB = db.New(testUtil.Postgres(), testUtil.Rabbit(), rabbit.PublishOptions{})

	// run the tests
	code := main.Run()

	// teardown
	testUtil.TearDown()
	os.Exit(code)
}

// TestMain will be called by `go test`, it can do setup/teardown/etc
// This is called for every file of tests not just once per invocation of `go test`
