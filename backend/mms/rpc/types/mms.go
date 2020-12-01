package types

import "time"

type MMS struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"account_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ProviderKey string    `json:"-"`
	MessageID   string    `json:"message_id"`
	MessageRef  string    `json:"message_ref"`
	Country     string    `json:"country"`
	Subject     string    `json:"subject"`
	Message     string    `json:"message"`
	ContentURLs []string  `json:"content_urls"`
	Recipient   string    `json:"recipient"`
	Sender      string    `json:"sender"`
	Status      string    `json:"status"`
	TrackLinks  bool      `json:"track_links"`
}

type SendParams struct {
	AccountID   string
	Subject     string
	Message     string
	Recipient   string
	Sender      string
	Country     string
	MessageRef  string
	ContentURLs []string
	TrackLinks  bool
	ProviderKey string
}

type SendReply struct {
	MMS *MMS
}

type UpdateStatusParams struct {
	ID          string
	MessageID   string
	Status      string
	Description string
}

type FindByIDParams struct {
	ID string
}

type FindByIDReply struct {
	MMS *MMS
}
