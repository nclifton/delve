package rpc

import (
	"time"

	"github.com/burstsms/mtmo-tp/backend/webhook/db"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/msg"
)

type PublishOptOutParams struct {
	Source        string         `json:"source"`
	Timestamp     time.Time      `json:"timestamp"`
	SourceMessage *SourceMessage `json:"source_message,omitempty"`
	AccountID     string         `json:"account_id"`
}

type SourceMessage struct {
	Type        string   `json:"type"`
	ID          string   `json:"id"`
	Recipient   string   `json:"recipient"`
	Sender      string   `json:"sender"`
	Message     string   `json:"message"`
	MessageRef  string   `json:"message_ref"`
	Subject     string   `json:"subject,omitempty"`
	ContentURLS []string `json:"content_urls,omitempty"`
}

func (s *Webhook) PublishOptOut(p PublishOptOutParams, r *NoReply) error {
	webhooks, err := s.db.FindByEvent(p.AccountID, "opt_out")
	if err != nil {
		return err
	}

	for _, w := range webhooks {
		err = s.db.Publish(msg.WebhookMessageSpec{
			URL:       w.URL,
			RateLimit: w.RateLimit,
			Payload:   msg.WebhookBody{Event: "opt_out", Data: p},
		}, db.RabbitPublishOptions{
			Exchange:     msg.WebhookMessage.Exchange,
			ExchangeType: msg.WebhookMessage.ExchangeType,
			RouteKey:     msg.WebhookMessage.RouteKey,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
