package types

import "time"

type OptOut struct {
	ID          string
	AccountID   string
	MessageID   string
	MessageType string
	Sender      string
	LinkID      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FindByLinkIDParams struct {
	LinkID string
}

type FindByLinkIDReply struct {
	*OptOut
}

type OptOutViaLinkParams struct {
	LinkID string
}

type OptOutViaLinkReply struct {
	*OptOut
}

type GenerateOptoutLinkParams struct {
	AccountID   string
	MessageID   string
	MessageType string
	Message     string
}

type GenerateOptoutLinkReply struct {
	Message string
}
