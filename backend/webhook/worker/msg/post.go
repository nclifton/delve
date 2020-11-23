package msg

var WebhookMessage = MessageSpec{
	Queue:        "webhook",
	Exchange:     "webhook",
	ExchangeType: "direct",
	RouteKey:     "",
}

type WebhookBody struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type WebhookMessageSpec struct {
	URL       string      `json:"url"`
	RateLimit int         `json:"rate_limit"`
	Payload   WebhookBody `json:"payload"`
}
