package msg

var MOMessage = MessageSpec{
	Queue:        "mo",
	Exchange:     "mo",
	ExchangeType: "direct",
	RouteKey:     "",
}

type MOMessageSpec struct {
	MessageID     string `json:"message_id"`
	Message       string `json:"message"`
	To            string `json:"to"`
	From          string `json:"from"`
	SARID         string `json:"sar_id"`
	SARPartNumber string `json:"sar_part_number"`
	SARParts      string `json:"sar_parts"`
}
