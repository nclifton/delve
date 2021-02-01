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
		`select id, account_id, address, mms_provider_key, channels, country, COALESCE(comment, '') as comment, created_at, updated_at
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

func (db *sqlDB) SenderFindByAccountId(ctx context.Context, accountId string) ([]Sender, error) {
	rows, err := db.sql.Query(
		ctx,
		`select id, account_id, address, mms_provider_key, channels, country, COALESCE(comment, ''), created_at, updated_at
		from sender
		where account_id = $1
		limit 100`,
		accountId,
	)
	if err != nil {
		return []Sender{}, err
	}
	defer rows.Close()

	ss := []Sender{}
	for rows.Next() {
		s := Sender{}
		var channels pgtype.EnumArray
		err := rows.Scan(
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
			return nil, err
		}
		err = channels.AssignTo(&s.Channels)
		if err != nil {
			return nil, err
		}
		ss = append(ss, s)
	}

	return ss, nil
}
