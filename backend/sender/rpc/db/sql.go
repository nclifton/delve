package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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
	var accountId sql.NullString
	var channels pgtype.EnumArray
	var mmsProviderKey sql.NullString
	var comment sql.NullString
	err := row.Scan(
		&s.ID,
		&accountId,
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
	s.AccountID = accountId.String
	err = channels.AssignTo(&s.Channels)
	if err != nil {
		return Sender{}, err
	}
	s.MMSProviderKey = mmsProviderKey.String
	s.Comment = comment.String
	return s, nil
}

func (db *sqlDB) InsertSenders(ctx context.Context, newSenders []Sender) ([]Sender, error) {

	insertSql := "insert into sender (account_id, address, channels, mms_provider_key, country, comment, created_at, updated_at)"
	returningSql := "returning " + senderSelect
	valuesSql := ""
	valuesRowsSql := make([]string, 0, len(newSenders))
	args := make([]interface{}, 0, len(newSenders)*6)
	idx := 1
	for _, newSender := range newSenders {
		valuesRowSql := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, localtimestamp, localtimestamp)", idx, idx+1, idx+2, idx+3, idx+4, idx+5)
		valuesRowsSql = append(valuesRowsSql, valuesRowSql)
		args = append(args, nilIfBlank(newSender.AccountID))
		args = append(args, nilIfBlank(newSender.Address))
		args = append(args, newSender.Channels)
		args = append(args, nilIfBlank(newSender.MMSProviderKey))
		args = append(args, nilIfBlank(newSender.Country))
		args = append(args, nilIfBlank(newSender.Comment))
		idx = idx + 6
	}
	valuesSql = strings.Join(valuesRowsSql, ",\n")
	sqlStr := fmt.Sprintf("%s\nVALUES\n%s\n%s", insertSql, valuesSql, returningSql)
	rows, err := db.sql.Query(ctx, sqlStr, args...)
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
	if len(newSenders) != len(ss) {
		err = rows.Err()
	} else {
		err = nil
	}

	return ss, err
}

func nilIfBlank(v string) interface{} {
	if len(v) == 0 {
		return nil
	}
	return v
}

func (db *sqlDB) SenderAddressExists(ctx context.Context, address string) (bool, error) {

	var exists bool
	err := db.sql.QueryRow(ctx, `SELECT exists (SELECT FROM sender WHERE address = $1)`, address).Scan(&exists)
	if err != nil && err != pgx.ErrNoRows {
		return false, err
	}

	return exists, nil
}

var senderEnums SenderEnums

func (db *sqlDB) GetSenderEnums(ctx context.Context) (SenderEnums, error) {

	//we should cache this locally so we don't need to make this query more than once during the lifetime of the service
	if senderEnums != nil {
		return senderEnums, nil
	}

	rows, err := db.sql.Query(ctx,
		`SELECT type.typname, array_agg(enum.enumlabel) as value from pg_enum as enum join pg_type as type on (type.oid = enum.enumtypid) group by type.typname`,
	)
	if err != nil {
		return SenderEnums{}, err
	}
	defer rows.Close()

	names := make([]string, 0)
	valueArrays := make([][]string, 0)
	var idx int = 0
	var pgValues pgtype.EnumArray
	enums := SenderEnums{}
	for rows.Next() {

		names = append(names, "")
		valueArrays = append(valueArrays, []string{})

		if err := rows.Scan(&names[idx], &pgValues); err != nil {
			return SenderEnums{}, err
		}

		if err := pgValues.AssignTo(&valueArrays[idx]); err != nil {
			return SenderEnums{}, err
		}

		enums[names[idx]] = valueArrays[idx]
		idx++
	}

	senderEnums = enums

	return enums, nil
}
