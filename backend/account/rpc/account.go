package rpc

import (
	"time"
)

type Account struct {
	ID        string    `bson:"_id" json:"_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	Name      string    `bson:"name" json:"name"`

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
	account, err := s.db.FindByAPIKey(p.Key)
	if err != nil {
		return err
	}
	r.Account = account
	return nil
}
