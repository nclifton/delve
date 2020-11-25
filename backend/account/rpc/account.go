package rpc

import (
	"time"
)

type Account struct {
	ID             string    `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Name           string    `json:"name"`
	Sender         []string  `json:"sender"`
	AlarisUsername string    `json:"alaris_username"`
	AlarisPassword string    `json:"alaris_password"`
	AlarisURL      string    `json:"alaris_url"`

	APIKeys []AccountAPIKey `json:"api_keys" bson:"api_keys"`
}

type AccountAPIKey struct {
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
	Description string    `json:"description" bson:"description"`
	Key         string    `json:"key" bson:"key"`
}

type FindByAPIKeyParams struct {
	Key string
}

type FindByAPIKeyReply struct {
	Account *Account
}

func (s *AccountService) FindByAPIKey(p FindByAPIKeyParams, r *FindByAPIKeyReply) error {

	var account *Account

	err := s.db.redis.Cached(
		"Account.FindByAPIKey:"+p.Key,
		&account,
		time.Minute*5,
		func() (interface{}, error) {
			return s.db.FindByAPIKey(p.Key)
		},
	)
	if err != nil {
		return err
	}
	r.Account = account
	return nil
}

type FindBySenderParams struct {
	Sender string
}

type FindBySenderReply struct {
	Account *Account
}

func (s *AccountService) FindBySender(p FindBySenderParams, r *FindBySenderReply) error {
	var account *Account

	err := s.db.redis.Cached(
		"Account.FindBySender:"+p.Sender,
		&account,
		time.Minute*5,
		func() (interface{}, error) {
			return s.db.FindBySender(p.Sender)
		},
	)
	if err != nil {
		return err
	}
	r.Account = account
	return nil
}
