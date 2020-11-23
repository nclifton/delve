package worker

type WebhookMessage struct {
	URL       string `json:"url"`
	RateLimit int    `json:"rate_limit"`
	Payload   []byte `json:"payload"`
}
