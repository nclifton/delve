package util

import (
	"context"
	"log"

	"github.com/NeowayLabs/wabbit"
	"github.com/NeowayLabs/wabbit/amqptest"
	amqptestServer "github.com/NeowayLabs/wabbit/amqptest/server"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TestUtil struct {
	postgres      *pgxpool.Pool
	rabbit        wabbit.Conn
	rabbitServer  *amqptestServer.AMQPServer
	migrate       *migrate.Migrate
	postgresURL   string
	migrationsURL string
}

func (t *TestUtil) resetMigration() {
	// initialise migrations
	m, err := migrate.New(
		"file://"+t.migrationsURL,
		t.postgresURL,
	)
	if err != nil {
		log.Fatalf("Failed to initialise golang-migrate connection: %s\n", err)
	}
	t.migrate = m
}

func (t *TestUtil) Setup() {
	t.resetMigration()

	// clear db
	if err := t.migrate.Drop(); err != nil {
		log.Fatalf("Failed to drop db: %s\n", err)
	}

	t.resetMigration()

	// apply migrations
	if err := t.migrate.Up(); err != nil {
		log.Fatalf("Failed to run migrations: %s\n", err)
	}
}

func (t *TestUtil) TearDown() {
	t.postgres.Close()
	err := t.migrate.Drop()
	if err != nil {
		log.Fatalf("Failed to drop test db: %s\n", err)
	}
	err = t.rabbitServer.Stop()
	if err != nil {
		log.Fatalf("Failed to stop test rabbit server: %s\n", err)
	}
}

func (t *TestUtil) Rabbit() wabbit.Conn {
	return t.rabbit
}

func (t *TestUtil) Postgres() *pgxpool.Pool {
	return t.postgres
}

func NewTestUtil(postgresURL, migrationsURL string) *TestUtil {
	var err error

	testPostgres, err := pgxpool.Connect(context.Background(), postgresURL)
	if err != nil {
		log.Fatalf("Failed to initialise postgres conn: %s\n", err)
	}

	testRabbitServer := amqptestServer.NewServer("amqp://localhost:5672/%2f")
	err = testRabbitServer.Start()
	if err != nil {
		log.Fatalf("Failed to start fake rabbit server: %s\n", err)
	}
	testRabbit, err := amqptest.Dial("amqp://localhost:5672/%2f")
	if err != nil {
		log.Fatalf("Failed to initialise fake rabbit client: %s\n", err)
	}

	return &TestUtil{
		postgres:      testPostgres,
		rabbit:        testRabbit,
		rabbitServer:  testRabbitServer,
		postgresURL:   postgresURL,
		migrationsURL: migrationsURL,
	}
}
