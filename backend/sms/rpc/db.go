package rpc

import (
	"context"

	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	SMSTable = "sms"
)

type db struct {
	postgres *pgxpool.Pool
	rabbit   rabbit.Conn
}

// New db interface

func NewDB(postgresURL string, rabbitmq rabbit.Conn) (*db, error) {
	postgres, err := pgxpool.Connect(context.Background(), postgresURL)
	if err != nil {
		return nil, err
	}

	return &db{postgres: postgres, rabbit: rabbitmq}, nil
}

type CommandTag = pgconn.CommandTag

func (db *db) Exec(sql string, args ...interface{}) (CommandTag, error) {
	return db.postgres.Exec(bg(), sql, args...)
}

func bg() context.Context {
	return context.Background()
}

type RabbitPublishOptions = rabbit.PublishOptions

func (db *db) Publish(msg interface{}, opts RabbitPublishOptions) error {
	return rabbit.Publish(db.rabbit, opts, msg)
}
