package rpc

import (
	"context"
	"database/sql"

	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/mms/rpc/types"
	"github.com/jackc/pgx/v4/pgxpool"
)

type db struct {
	postgres *pgxpool.Pool
	rabbit   rabbit.Conn
	opts     RabbitPublishOptions
}

// New db interface
func NewDB(postgresURL string, rabbitmq rabbit.Conn, opts RabbitPublishOptions) (*db, error) {
	postgres, err := pgxpool.Connect(context.Background(), postgresURL)
	if err != nil {
		return nil, err
	}

	return &db{postgres: postgres, rabbit: rabbitmq, opts: opts}, nil
}

func (db *db) FindByID(ctx context.Context, id string) (*types.MMS, error) {
	var mms types.MMS

	query := `
		SELECT id, account_id, created_at, updated_at, provider_key, message_id, message_ref,
			country, subject, message, content_urls, recipient, sender, status,
			track_links
		FROM mms
		WHERE id = $1
	`

	row := db.postgres.QueryRow(ctx, query, id)

	var msgID sql.NullString

	err := row.Scan(
		&mms.ID,
		&mms.AccountID,
		&mms.CreatedAt,
		&mms.UpdatedAt,
		&mms.ProviderKey,
		&msgID,
		&mms.MessageRef,
		&mms.Country,
		&mms.Subject,
		&mms.Message,
		&mms.ContentURLs,
		&mms.Recipient,
		&mms.Sender,
		&mms.Status,
		&mms.TrackLinks,
	)
	if err != nil {
		return nil, err
	}

	mms.MessageID = msgID.String

	return &mms, nil
}

func (db *db) FindByIDAndAccountID(ctx context.Context, id, accountID string) (*types.MMS, error) {
	var mms types.MMS

	sql := `
		SELECT id, account_id, created_at, updated_at, provider_key, message_id, message_ref,
			country, subject, message, content_urls, recipient, sender, status,
			track_links
		FROM mms
		WHERE id = $1 and account_id = $2
	`

	row := db.postgres.QueryRow(ctx, sql, id, accountID)

	err := row.Scan(
		&mms.ID,
		&mms.AccountID,
		&mms.CreatedAt,
		&mms.UpdatedAt,
		&mms.ProviderKey,
		&mms.MessageID,
		&mms.MessageRef,
		&mms.Country,
		&mms.Subject,
		&mms.Message,
		&mms.ContentURLs,
		&mms.Recipient,
		&mms.Sender,
		&mms.Status,
		&mms.TrackLinks,
	)
	if err != nil {
		return nil, err
	}

	return &mms, nil
}

func (db *db) InsertMMS(ctx context.Context, mms types.MMS) (*types.MMS, error) {
	sql := `
		INSERT INTO mms (id, account_id, created_at, updated_at, provider_key,
			message_ref, country, subject, message,
			content_urls, recipient, sender, status, track_links)
		VALUES ($1, $2, NOW(), NOW(), $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`

	if err := db.postgres.QueryRow(ctx, sql,
		mms.ID,
		mms.AccountID,
		mms.ProviderKey,
		mms.MessageRef,
		mms.Country,
		mms.Subject,
		mms.Message,
		mms.ContentURLs,
		mms.Recipient,
		mms.Sender,
		mms.Status,
		mms.TrackLinks,
	).Scan(&mms.ID, &mms.CreatedAt, &mms.UpdatedAt); err != nil {
		return &types.MMS{}, err
	}

	return &mms, nil
}

func (db *db) UpdateStatus(ctx context.Context, id, messageID, status string) error {
	sql := `
		UPDATE mms
		SET status = $3, message_id = $2
		WHERE id = $1
	`

	_, err := db.postgres.Exec(ctx, sql, id, messageID, status)

	return err
}

func bg() context.Context {
	return context.Background()
}

type RabbitPublishOptions = rabbit.PublishOptions

func (db *db) Publish(msg interface{}, routeKey string) error {
	publishOpts := RabbitPublishOptions{
		Exchange:     db.opts.Exchange,
		ExchangeType: db.opts.ExchangeType,
	}
	return rabbit.Publish(db.rabbit, publishOpts, msg)
}
