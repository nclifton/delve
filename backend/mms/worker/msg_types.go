package worker

const (
	MMSSendRouteKey  = ""
	MMSSendQueueName = "mms.single"
)

type Job struct {
	ID          string   `json:"_id"`
	AccountID   string   `json:"account_id"`
	Sender      string   `json:"sender"`
	Subject     string   `json:"subject"`
	ContentURLs []string `json:"content_urls"`
	Recipient   string   `json:"recipient"`
	ProviderKey string   `json:"provider_key"`
	Message     string   `json:"message"`
}
