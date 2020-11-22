package msg

var SMSSendMessage = MessageSpec{
	Queue:        "send.single",
	Exchange:     "send.single",
	ExchangeType: "direct",
	RouteKey:     "",
}

type SMSSendMessageSpec struct {
	ID         string `json:"id"`
	AccountID  string `json:"account_id"`
	Message    string `json:"message"`
	Recipient  string `json:"recipient"`
	Sender     string `json:"sender"`
	AlarisUser string `json:"alaris_user"`
	AlarisPass string `json:"alaris_pass"`
}
