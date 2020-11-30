package types

import (
	"time"
)

type PublishOptOutParams struct {
	Source        string         `json:"source"`
	Timestamp     time.Time      `json:"timestamp"`
	SourceMessage *SourceMessage `json:"source_message,omitempty"`
	AccountID     string         `json:"account_id"`
}

type SourceMessage struct {
	Type        string   `json:"type"`
	ID          string   `json:"id"`
	Recipient   string   `json:"recipient"`
	Sender      string   `json:"sender"`
	Message     string   `json:"message"`
	MessageRef  string   `json:"message_ref"`
	Subject     string   `json:"subject,omitempty"`
	ContentURLS []string `json:"content_urls,omitempty"`
}
