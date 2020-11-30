package types

import "time"

const Name = "OptOut"

type NoParams struct{}
type NoReply struct{}

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

const FindByLinkID = "FindByLinkID"

type FindByLinkIDParams struct {
	LinkID string
}

type FindByLinkIDReply struct {
	*OptOut
}

const OptOutViaLink = "OptOutViaLink"

type OptOutViaLinkParams struct {
	LinkID string
}

type OptOutViaLinkReply struct {
	*OptOut
}

const GenerateOptOutLink = "GenerateOptOutLink"

type GenerateOptOutLinkParams struct {
	AccountID   string
	MessageID   string
	MessageType string
	Message     string
}

type GenerateOptOutLinkReply struct {
	Message string
}

const OptOutViaMsg = "OptOutViaMsg"

type OptOutViaMsgParams struct {
	Message     string
	AccountID   string
	MessageType string
	MessageID   string
}
