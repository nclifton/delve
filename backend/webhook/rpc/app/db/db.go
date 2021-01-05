package db

import (
	"context"
	"time"
)

type Webhook struct {
	ID        int64
	AccountID string
	Event     string
	Name      string
	URL       string
	RateLimit int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DB interface {
	InsertWebhook(ctx context.Context, accountID, event, name, url string, rateLimit int32) (Webhook, error)
	FindWebhook(ctx context.Context, accountID string) ([]Webhook, error)
	FindWebhookByEvent(ctx context.Context, accountID string, event string) ([]Webhook, error)
	DeleteWebhook(ctx context.Context, id int64, accountID string) error
	UpdateWebhook(ctx context.Context, id int64, accountID, event, name, url string, rateLimit int32) (Webhook, error)
}
