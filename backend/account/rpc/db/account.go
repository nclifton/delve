package db

import (
	"context"

	"github.com/burstsms/mtmo-tp/backend/lib/errorlib"
	"github.com/jackc/pgx/v4"
)

const msgAccountNotFound = "account not found"

const accountFields = `a.id, 
a.created_at, 
a.updated_at, 
a.name, 
a.alaris_username, 
a.alaris_password, 
a.alaris_url`

func scanAccount(row pgx.Row, account *Account) error {
	return row.Scan(
		&account.ID,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.Name,
		&account.AlarisUsername,
		&account.AlarisPassword,
		&account.AlarisURL,
	)
}

func (db *sqlDB) FindAccountByAPIKey(ctx context.Context, key string) (Account, error) {
	account := Account{}

	sql := `
		SELECT ` + accountFields + `
		FROM account a
		LEFT JOIN account_api_keys as ak ON a.id = ak.account_id
		WHERE ak.key = $1;
	`

	row := db.sql.QueryRow(ctx, sql, key)
	if err := scanAccount(row, &account); err != nil {
		if err == pgx.ErrNoRows {
			return Account{}, errorlib.NotFoundErr{Message: msgAccountNotFound}
		}

		return Account{}, err
	}

	return account, nil
}

func (db *sqlDB) FindAccountByID(ctx context.Context, id string) (Account, error) {
	account := Account{}

	query := `
		SELECT ` + accountFields + `
		FROM account a
		WHERE id = $1;
	`

	row := db.sql.QueryRow(ctx, query, id)
	if err := scanAccount(row, &account); err != nil {
		if err == pgx.ErrNoRows {
			return Account{}, errorlib.NotFoundErr{Message: msgAccountNotFound}
		}

		return Account{}, err
	}

	return account, nil
}
