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
	opts     RabbitPublishOptions
}

// New db interface
func New(postgres *pgxpool.Pool, rabbitmq rabbit.Conn, opts RabbitPublishOptions) *DB {
	return &DB{postgres: postgres, rabbit: rabbitmq, opts: opts}
}

type CommandTag = pgconn.CommandTag

func (db *DB) Exec(sql string, args ...interface{}) (CommandTag, error) {
	return db.postgres.Exec(bg(), sql, args...)
}

func bg() context.Context {
	return context.Background()
}

type RabbitPublishOptions = rabbit.PublishOptions

func (db *DB) Publish(msg interface{}) error {
	return rabbit.Publish(db.rabbit, db.opts, msg)
}
