package types

import "time"

type PublishLinkHitParams struct {
	URL           string
	Hits          int
	Timestamp     time.Time
	SourceMessage *SourceMessage
	AccountID     string
}
