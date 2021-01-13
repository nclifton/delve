package db

import (
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

type DB interface{}
