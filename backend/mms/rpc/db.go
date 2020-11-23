package rpc

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type db struct {
	postgres *pgxpool.Pool
}

// New db interface
func NewDB(postgresURL string) (*db, error) {
	postgres, err := pgxpool.Connect(context.Background(), postgresURL)
	if err != nil {
		return nil, err
	}

	return &db{postgres: postgres}, nil
}

func (db *db) FindByID(ctx context.Context, id, accountID string) (*MMS, error) {
	var mms MMS

	sql := `
		SELECT id, created_at, updated_at, provider_key, message_id, message_ref,
			country, subject, messsage, content_urls, recipient, sender, status,
			shorten_urls, unsub
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
		&mms.ShortenURLs,
		&mms.Unsub,
	)
	if err != nil {
		return nil, err
	}

	return &mms, nil
}

func (db *db) InsertMMS(ctx context.Context, mms MMS) (*MMS, error) {
	sql := `
		INSERT INTO mms (account_id, created_at, updated_at, provider_key,
			message_ref, country, subject, message,
			content_urls, recipient, sender, status, shorten_urls)
		VALUES ($1, NOW(), NOW(), $2, $3, $4, $5, $6, $7, $8, $9, $10) 
		RETURNING id, created_at, updated_at
	`

	if err := db.postgres.QueryRow(ctx, sql,
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
		mms.ShortenURLs,
	).Scan(&mms.ID, &mms.CreatedAt, &mms.UpdatedAt); err != nil {
		return &MMS{}, err
	}

	return &mms, nil
}
