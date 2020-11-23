package db

import (
	"context"

	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

// define table names??
// TODO discovery on how to do postgres migrations and models
const (
	WebhookTable = "webhook"
)

// DB will wrap the underlying connection so that exported methods
// can concistently drive database operations
type DB struct {
	postgres *pgxpool.Pool
	rabbit   rabbit.Conn
}

// New db interface
func New(postgres *pgxpool.Pool, rabbitmq rabbit.Conn) *DB {
	return &DB{postgres: postgres, rabbit: rabbitmq}
}

type CommandTag = pgconn.CommandTag

func (db *DB) Exec(sql string, args ...interface{}) (CommandTag, error) {
	return db.postgres.Exec(bg(), sql, args...)
}

func bg() context.Context {
	return context.Background()
}

type RabbitPublishOptions = rabbit.PublishOptions

func (db *DB) Publish(msg interface{}, opts RabbitPublishOptions) error {
	return rabbit.Publish(db.rabbit, opts, msg)
}
