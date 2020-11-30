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

func (db *db) InsertTrackLink(ctx context.Context, accountID, messageID, messageType, url string) (*TrackLink, error) {
	tracklink := TrackLink{}
	err := db.postgres.QueryRow(
		ctx,
		`INSERT INTO track_link(account_id, message_id, message_type, created_at, updated_at, url)
		VALUES($1, $2, $3, now(), now(), $4)
		RETURNING id, account_id, message_id, message_type, track_link_id, created_at, updated_at, url, hits`,
		accountID,
		messageID,
		messageType,
		url,
	).Scan(
		&tracklink.ID,
		&tracklink.AccountID,
		&tracklink.MessageID,
		&tracklink.MessageType,
		&tracklink.TrackLinkID,
		&tracklink.CreatedAt,
		&tracklink.UpdatedAt,
		&tracklink.URL,
		&tracklink.Hits,
	)
	return &tracklink, err
}

func (db *db) FindTrackLinkByTrackLinkID(ctx context.Context, accountID, tracklinkID string) (*TrackLink, error) {
	tracklink := TrackLink{}
	err := db.postgres.QueryRow(
		ctx,
		`SELECT id, account_id, message_id, message_type, track_link_id, created_at, updated_at, url, hits
		FROM track_link
		WHERE account_id = $1 AND track_link_id = $2`,
		accountID,
		tracklinkID,
	).Scan(
		&tracklink.ID,
		&tracklink.AccountID,
		&tracklink.MessageID,
		&tracklink.MessageType,
		&tracklink.TrackLinkID,
		&tracklink.CreatedAt,
		&tracklink.UpdatedAt,
		&tracklink.URL,
		&tracklink.Hits,
	)
	return &tracklink, err
}

func (db *db) IncrementTrackLinkHits(ctx context.Context, accountID, tracklinkID string) (*TrackLink, error) {
	tracklink := TrackLink{}
	err := db.postgres.QueryRow(
		ctx,
		`UPDATE track_link
		SET hits = hits + 1
		WHERE account_id = $1 AND track_link_id = $2
		RETURNING id, account_id, message_id, message_type, track_link_id, created_at, updated_at, url, hits`,
		accountID,
		tracklinkID,
	).Scan(
		&tracklink.ID,
		&tracklink.AccountID,
		&tracklink.MessageID,
		&tracklink.MessageType,
		&tracklink.TrackLinkID,
		&tracklink.CreatedAt,
		&tracklink.UpdatedAt,
		&tracklink.URL,
		&tracklink.Hits,
	)
	return &tracklink, err
}
