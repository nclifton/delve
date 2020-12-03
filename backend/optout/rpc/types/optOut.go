package types

import (
	"fmt"
	"time"
)

const Name = "OptOut"

type NoParams struct{}
type NoReply struct{}

type OptOut struct {
	ID          string
	AccountID   string
	MessageID   string
	MessageType string
	Sender      string // TODO: need to be added to the DB?
	LinkID      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func OptOutEqual(expected, actual OptOut) bool {
	if actual.ID == "" {
		fmt.Println("detected OptOut id empty")
		return false
	}

	if expected.AccountID != actual.AccountID {
		fmt.Println("detected OptOut not equal due to change in AccountID")
		return false
	}

	if expected.MessageID != actual.MessageID {
		fmt.Println("detected OptOut not equal due to change in MessageID")
		return false
	}

	if expected.MessageType != actual.MessageType {
		fmt.Println("detected OptOut not equal due to change in MessageType")
		return false
	}

	if actual.LinkID == "" {
		fmt.Println("detected OptOut LinkID empty")
		return false
	}

	return true
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

type GenerateOptOutLinkParams struct {
	AccountID   string
	MessageID   string
	MessageType string
	Message     string
}

type GenerateOptOutLinkReply struct {
	Message string
}

type OptOutViaMsgParams struct {
	Message     string
	AccountID   string
	MessageType string
	MessageID   string
}
