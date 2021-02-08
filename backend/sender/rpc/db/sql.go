package db

import (
	"context"
	"database/sql"
	"fmt"

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

func (db *sqlDB) FindSenderByAddressAndAccountID(ctx context.Context, accountId, address string) (Sender, error) {

	row := db.sql.QueryRow(ctx, fmt.Sprintf("select %s from sender where account_id = $1 and address = $2", senderSelect), accountId, address)
	s, err := scanSenderRow(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return s, errorlib.NotFoundErr{Message: "sender not found"}
		}
		return s, err
	}

	return s, nil
}
func (db *sqlDB) FindSendersByAccountId(ctx context.Context, accountId string) ([]Sender, error) {

	rows, err := db.sql.Query(ctx, fmt.Sprintf("select %s from sender where account_id = $1 limit 100", senderSelect), accountId)
	if err != nil {
		return []Sender{}, err
	}
	defer rows.Close()

	ss := []Sender{}
	for rows.Next() {
		s, err := scanSenderRow(rows)
		if err != nil {
			return []Sender{}, err
		}
		ss = append(ss, s)
	}

	return ss, nil
}

func (db *sqlDB) FindSendersByAddress(ctx context.Context, address string) ([]Sender, error) {

	rows, err := db.sql.Query(ctx, fmt.Sprintf("select %s from sender where address = $1", senderSelect), address)
	if err != nil {
		return []Sender{}, err
	}
	defer rows.Close()

	ss := []Sender{}
	for rows.Next() {
		s, err := scanSenderRow(rows)
		if err != nil {
			return []Sender{}, err
		}
		ss = append(ss, s)
	}

	return ss, nil
}

const senderSelect string = "id, account_id, address, mms_provider_key, channels, country, comment, created_at, updated_at"

func scanSenderRow(row pgx.Row) (Sender, error) {
	s := Sender{}
	var channels pgtype.EnumArray
	var mmsProviderKey sql.NullString
	var comment sql.NullString
	err := row.Scan(
		&s.ID,
		&s.AccountID,
		&s.Address,
		&mmsProviderKey,
		&channels,
		&s.Country,
		&comment,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		return Sender{}, err
	}
	err = channels.AssignTo(&s.Channels)
	if err != nil {
		return Sender{}, err
	}
	s.MMSProviderKey = mmsProviderKey.String
	s.Comment = comment.String
	return s, nil
}
