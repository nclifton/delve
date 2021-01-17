package db

import (
	"context"
	"time"
)

type Sender struct {
	ID             string
	AccountID      string
	Address        string
	MMSProviderKey string
	Channels       []string
	Country        string
	Comment        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type DB interface {
	SenderFindByAddress(ctx context.Context, accountId, address string) (Sender, error)
	SenderFindByAccountId(ctx context.Context, accountId string) ([]Sender, error)
}
