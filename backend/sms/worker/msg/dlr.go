package msg

import "time"

var DLRMessage = MessageSpec{
	Queue:        "dlr",
	Exchange:     "dlr",
	ExchangeType: "direct",
	RouteKey:     "",
}

type DLRMessageSpec struct {
	MessageID  string    `json:"message_id"`
	State      string    `json:"state"`
	ReasonCode string    `json:"reason_code"`
	To         string    `json:"to"`
	Time       time.Time `json:"time"`
	MCC        string    `json:"mcc"`
	MNC        string    `json:"mnc"`
}
