package types

import "time"

type SMS struct {
	ID         string    `json:"id"`
	AccountID  string    `json:"account_id"`
	MessageID  string    `json:"message_id"`
	Recipient  string    `json:"recipient"`
	Sender     string    `json:"sender"`
	Country    string    `json:"country"`
	MessageRef string    `json:"message_ref"`
	Message    string    `json:"message"`
	Status     string    `json:"status"`
	SMSCount   int       `json:"sms_count"`
	GSM        bool      `json:"gsm"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	TrackLinks bool      `json:"track_links"`
}

type SendParams struct {
	AccountID  string
	Message    string
	Recipient  string
	Sender     string
	Country    string
	MessageRef string
	AlarisUser string
	AlarisPass string
	AlarisURL  string
	TrackLinks bool
}

type SendReply struct {
	SMS *SMS
}
type MarkSentParams struct {
	ID        string
	AccountID string
	MessageID string
}

type MarkFailedParams struct {
	AccountID string
	ID        string
}

type FindByIDParams struct {
	AccountID string
	ID        string
}

type FindByIDReply struct {
	*SMS
}
