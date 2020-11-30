package types

import "time"

type QueueDLRParams struct {
	MessageID  string
	State      string
	ReasonCode string
	To         string
	Time       time.Time
	MCC        string
	MNC        string
}

type ProcessDLRParams struct {
	MessageID  string
	State      string
	ReasonCode string
	To         string
	Time       time.Time
	MCC        string
	MNC        string
}
