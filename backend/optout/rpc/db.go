package rpc

import (
	"context"
	"fmt"

	"github.com/burstsms/mtmo-tp/backend/lib/errors"
	types "github.com/burstsms/mtmo-tp/backend/optout/rpc/types"
	"github.com/jackc/pgx/v4"
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

func (db *db) FindOptOutByLinkID(ctx context.Context, linkID string) (*types.OptOut, error) {
	optOut := types.OptOut{}

	if err := db.postgres.QueryRow(
		ctx,
		`SELECT id, account_id, message_id, message_type, link_id, created_at, updated_at
		FROM opt_out
		WHERE link_id = $1`,
		linkID,
	).Scan(
		&optOut.ID,
		&optOut.AccountID,
		&optOut.MessageID,
		&optOut.MessageType,
		&optOut.LinkID,
		&optOut.CreatedAt,
		&optOut.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return &types.OptOut{}, errors.NotFoundErr{Err: fmt.Errorf("optOut not found")}
		}
		return &types.OptOut{}, err
	}

	return &optOut, nil
}

func (db *db) InsertOptOut(ctx context.Context, accountID, messageID, messageType string) (*types.OptOut, error) {
	optOut := types.OptOut{}

	err := db.postgres.QueryRow(
		ctx,
		`INSERT INTO opt_out(account_id, message_id, message_type, created_at, updated_at)
		VALUES($1, $2, $3, now(), now())
		RETURNING id, account_id, message_id, message_type, link_id, created_at, updated_at`,
		accountID,
		messageID,
		messageType,
	).Scan(
		&optOut.ID,
		&optOut.AccountID,
		&optOut.MessageID,
		&optOut.MessageType,
		&optOut.LinkID,
		&optOut.CreatedAt,
		&optOut.UpdatedAt,
	)

	return &optOut, err
}
