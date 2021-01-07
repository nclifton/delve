package types

import "time"

type Account struct {
	ID             string    `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Name           string    `json:"name"`
	SenderSMS      []string  `json:"sender_sms"`
	SenderMMS      []string  `json:"sender_mms"`
	AlarisUsername string    `json:"alaris_username"`
	AlarisPassword string    `json:"alaris_password"`
	AlarisURL      string    `json:"alaris_url"`
	MMSProviderKey string    `json:"mms_provider_key"`

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

type FindBySenderParams struct {
	Sender string
}

type FindBySenderReply struct {
	Account *Account
}

type FindByIDParams struct {
	ID string
}

type FindByIDReply struct {
	Account *Account
}
