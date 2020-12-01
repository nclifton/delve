package types

import "time"

type SMS struct {
	ID         string
	AccountID  string
	MessageID  string
	Recipient  string
	Sender     string
	Country    string
	MessageRef string
	Message    string
	Status     string
	SMSCount   int
	GSM        bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	TrackLinks bool
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
