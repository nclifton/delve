package types

import (
	"time"
)

type PublishSMSStatusUpdateParams struct {
	AccountID       string    `json:"account_id"`
	SMSID           string    `json:"sms_id"`
	MessageRef      string    `json:"message_ref"`
	Recipient       string    `json:"recipient"`
	Sender          string    `json:"sender"`
	Status          string    `json:"status"`
	StatusUpdatedAt time.Time `json:"status_updated_at"`
}

type PublishMOParams struct {
	AccountID   string       `json:"account_id"`
	SMSID       string       `json:"sms_id"`
	Message     string       `json:"message"`
	Recipient   string       `json:"recipient"`
	Sender      string       `json:"sender"`
	ReceivedAt  time.Time    `json:"received_at"`
	LastMessage *LastMessage `json:"last_message,omitempty"`
}

type LastMessage struct {
	Type        string   `json:"type"`
	ID          string   `json:"id"`
	Recipient   string   `json:"recipient"`
	Sender      string   `json:"sender"`
	Message     string   `json:"message"`
	MessageRef  string   `json:"message_ref"`
	Subject     string   `json:"subject,omitempty"`
	ContentURLS []string `json:"content_urls,omitempty"`
}
