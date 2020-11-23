package msg

var WebhookMessage = MessageSpec{
	Queue:        "webhook",
	Exchange:     "webhook",
	ExchangeType: "direct",
	RouteKey:     "",
}

type WebhookMessageSpec struct {
	URL       string `json:"url"`
	RateLimit int    `json:"rate_limit"`
	Payload   []byte `json:"payload"`
}
