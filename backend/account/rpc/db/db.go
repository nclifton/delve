package db

import (
	"context"
	"time"
)

type Account struct {
	ID             string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Name           string
	AlarisUsername string
	AlarisPassword string
	AlarisURL      string
}

type AccountAPIKey struct {
	ID          string
	AccountID   string
	Description string
	Key         string
}

type DB interface {
	FindAccountByAPIKey(ctx context.Context, key string) (Account, error)
	FindAccountByID(ctx context.Context, id string) (Account, error)
}
