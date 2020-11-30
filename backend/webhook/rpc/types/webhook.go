package types

import (
	"time"
)

type WebhookRecord struct {
	ID        int64     `json:"id"`
	AccountID string    `json:"account_id"`
	Event     string    `json:"event"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	RateLimit int       `json:"rate_limit"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FindParams struct {
	AccountID string
}

type FindReply struct {
	Webhooks []WebhookRecord `json:"webhooks"`
}

type InsertParams struct {
	AccountID string `json:"account_id"`
	Event     string `json:"event"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	RateLimit int    `json:"rate_limit"`
}

type InsertReply struct {
	Webhook WebhookRecord
}

type DeleteParams struct {
	AccountID string `json:"account_id"`
	ID        string `json:"ids"`
}

type UpdateParams struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	Event     string `json:"event"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	RateLimit int    `json:"rate_limit"`
}

type UpdateReply struct {
	Webhook WebhookRecord
}
