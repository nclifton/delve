package db

import (
	"context"
	"errors"
	"time"

	"github.com/burstsms/mtmo-tp/backend/lib/errorlib"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type sqlDB struct {
	sql *pgxpool.Pool
}

func NewSQLDB(db *pgxpool.Pool) DB {
	return &sqlDB{
		sql: db,
	}
}

func (db *sqlDB) InsertWebhook(ctx context.Context, accountID, event, name, url string, rateLimit int32) (Webhook, error) {
	row := db.sql.QueryRow(
		ctx,
		`insert into webhook (account_id, event, name, url, rate_limit, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7)
		returning id, account_id, event, name, url, rate_limit, created_at, updated_at`,
		accountID,
		event,
		name,
		url,
		rateLimit,
		time.Now(),
		time.Now(),
	)

	wh := Webhook{}
	err := row.Scan(
		&wh.ID,
		&wh.AccountID,
		&wh.Event,
		&wh.Name,
		&wh.URL,
		&wh.RateLimit,
		&wh.CreatedAt,
		&wh.UpdatedAt,
	)

	return wh, err
}

func (db *sqlDB) FindWebhook(ctx context.Context, accountID string) ([]Webhook, error) {
	rows, err := db.sql.Query(
		ctx,
		`select id, account_id, event, name, url, rate_limit, created_at, updated_at
		from webhook where account_id = $1
		order by updated_at desc
		limit 100`,
		accountID,
	)
	if err != nil {
		return []Webhook{}, err
	}
	defer rows.Close()

	whs := []Webhook{}
	for rows.Next() {
		wh := Webhook{}
		err := rows.Scan(
			&wh.ID,
			&wh.AccountID,
			&wh.Event,
			&wh.Name,
			&wh.URL,
			&wh.RateLimit,
			&wh.CreatedAt,
			&wh.UpdatedAt,
		)
		if err != nil {
			return whs, err
		}

		whs = append(whs, wh)
	}

	return whs, nil
}

func (db *sqlDB) FindWebhookByEvent(ctx context.Context, accountID string, event string) ([]Webhook, error) {
	rows, err := db.sql.Query(
		ctx,
		`select id, account_id, event, name, url, rate_limit, created_at, updated_at
		from webhook
		where account_id = $1 and event = $2
		order by updated_at desc
		limit 100`,
		accountID,
		event,
	)
	if err != nil {
		return []Webhook{}, err
	}
	defer rows.Close()

	whs := []Webhook{}
	for rows.Next() {
		wh := Webhook{}
		err := rows.Scan(
			&wh.ID,
			&wh.AccountID,
			&wh.Event,
			&wh.Name,
			&wh.URL,
			&wh.RateLimit,
			&wh.CreatedAt,
			&wh.UpdatedAt,
		)
		if err != nil {
			return whs, err
		}

		whs = append(whs, wh)
	}

	return whs, nil
}

func (db *sqlDB) FindWebhookByID(ctx context.Context, accountID string, webhookID string) (Webhook, error) {
	row := db.sql.QueryRow(
		ctx,
		`select id, account_id, event, name, url, rate_limit, created_at, updated_at
		from webhook
		where account_id = $1 and id = $2`,
		accountID,
		webhookID,
	)

	wh := Webhook{}
	err := row.Scan(
		&wh.ID,
		&wh.AccountID,
		&wh.Event,
		&wh.Name,
		&wh.URL,
		&wh.RateLimit,
		&wh.CreatedAt,
		&wh.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return wh, errorlib.NotFoundErr{Message: "webhook not found"}
		}
		return wh, err
	}

	return wh, nil
}

func (db *sqlDB) DeleteWebhook(ctx context.Context, id int64, accountID string) error {

	ct, err := db.sql.Exec(
		ctx,
		`delete from webhook where
		id = $1 and account_id = $2`,
		id,
		accountID,
	)
	if err != nil {
		return err
	}

	count := ct.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("not found")
	}

	return nil
}

func (db *sqlDB) UpdateWebhook(ctx context.Context, id int64, accountID, event, name, url string, rateLimit int32) (Webhook, error) {
	row := db.sql.QueryRow(
		ctx,
		`update webhook set
		event = $1,
		name = $2,
		url = $3,
		rate_limit = $4,
		updated_at = $5
		where id = $6 and account_id = $7
		returning id, account_id, event, name, url, rate_limit, created_at, updated_at`,
		event,
		name,
		url,
		rateLimit,
		time.Now(),
		id,
		accountID,
	)

	wh := Webhook{}
	err := row.Scan(
		&wh.ID,
		&wh.AccountID,
		&wh.Event,
		&wh.Name,
		&wh.URL,
		&wh.RateLimit,
		&wh.CreatedAt,
		&wh.UpdatedAt,
	)
	if err != nil && err == pgx.ErrNoRows {
		return wh, errors.New("not found")
	}
	if err != nil {
		return wh, err
	}

	return wh, nil
}
