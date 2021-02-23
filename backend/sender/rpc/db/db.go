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

type SenderEnums map[string][]string

type DB interface {
	FindSenderByAddressAndAccountID(ctx context.Context, accountId, address string) (Sender, error)
	FindSendersByAccountId(ctx context.Context, accountId string) ([]Sender, error)
	FindSendersByAddress(ctx context.Context, address string) ([]Sender, error)
	InsertSenders(ctx context.Context, senders []Sender) ([]Sender, error)
	SenderAddressExists(ctx context.Context, address string) (bool, error)
	GetSenderEnums(ctx context.Context) (SenderEnums, error)
}
