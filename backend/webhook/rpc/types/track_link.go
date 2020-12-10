package types

import "time"

type PublishLinkHitParams struct {
	URL           string         `json:"url"`
	Hits          int            `json:"hits"`
	Timestamp     time.Time      `json:"timestamp"`
	SourceMessage *SourceMessage `json:"source_message,omitempty"`
	AccountID     string         `json:"account_id"`
}
