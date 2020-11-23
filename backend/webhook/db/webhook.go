package db

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
)

type webhook struct {
	ID        int64
	AccountID string
	Event     string
	Name      string
	URL       string
	RateLimit int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (db *DB) Insert(accountID, event, name, url string, rateLimit int) (webhook, error) {
	row := db.postgres.QueryRow(
		bg(),
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

	wh := webhook{}
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

func (db *DB) Find(accountID string) ([]webhook, error) {
	rows, err := db.postgres.Query(
		bg(),
		`select id, account_id, event, name, url, rate_limit, created_at, updated_at
		from webhook where account_id = $1
		order by updated_at desc
		limit 100`,
		accountID,
	)
	if err != nil {
		return []webhook{}, err
	}
	defer rows.Close()

	whs := []webhook{}
	for rows.Next() {
		wh := webhook{}
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

func (db *DB) FindByEvent(accountID string, event string) ([]webhook, error) {
	rows, err := db.postgres.Query(
		bg(),
		`select id, account_id, event, name, url, rate_limit, created_at, updated_at
		from webhook
		where account_id = $1 and event = $2
		order by updated_at desc
		limit 100`,
		accountID,
		event,
	)
	if err != nil {
		return []webhook{}, err
	}
	defer rows.Close()

	whs := []webhook{}
	for rows.Next() {
		wh := webhook{}
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

func (db *DB) Delete(accountID, id string) error {

	ct, err := db.Exec(
		`delete from webhook where
		account_id = $1 and id = $2`,
		accountID,
		id,
	)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return errors.New("not found")
	}

	return nil
}

func (db *DB) Update(webhookID, accountID, event, name, url string, rateLimit int) (webhook, error) {
	row := db.postgres.QueryRow(
		bg(),
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
		webhookID,
		accountID,
	)

	wh := webhook{}
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
