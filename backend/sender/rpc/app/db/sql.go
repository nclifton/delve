package db

import (
	"context"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/burstsms/mtmo-tp/backend/lib/errorlib"
)

type sqlDB struct {
	sql *pgxpool.Pool
}

func NewSQLDB(db *pgxpool.Pool) DB {
	return &sqlDB{
		sql: db,
	}
}

func (db *sqlDB) SenderFindByAddress(ctx context.Context, accountId, address string) (Sender, error) {
	row := db.sql.QueryRow(
		ctx,
		`select id, account_id, address, mms_provider_key, channels, country, comment, created_at, updated_at
		from sender
		where account_id = $1 and address = $2`,
		accountId,
		address,
	)
	s := Sender{}
	var channels pgtype.EnumArray
	err := row.Scan(
		&s.ID,
		&s.AccountID,
		&s.Address,
		&s.MMSProviderKey,
		&channels,
		&s.Country,
		&s.Comment,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return s, errorlib.NotFoundErr{Message: "sender not found"}
		}
		return s, err
	}
	err = channels.AssignTo(&s.Channels)
	if err != nil {
		return s, err
	}
	return s, nil
}
