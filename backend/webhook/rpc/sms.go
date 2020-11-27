package rpc

import (
	"time"

	"github.com/burstsms/mtmo-tp/backend/webhook/db"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/msg"
)

type PublishSMSStatusUpdateParams struct {
	AccountID       string    `json:"account_id"`
	SMSID           string    `json:"sms_id"`
	MessageRef      string    `json:"message_ref"`
	Recipient       string    `json:"recipient"`
	Sender          string    `json:"sender"`
	Status          string    `json:"status"`
	StatusUpdatedAt time.Time `json:"status_updated_at"`
}

func (s *Webhook) PublishSMSStatusUpdate(p PublishSMSStatusUpdateParams, r *NoReply) error {
	webhooks, err := s.db.FindByEvent(p.AccountID, "sms_status")
	if err != nil {
		return err
	}

	for _, w := range webhooks {
		err = s.db.Publish(msg.WebhookMessageSpec{
			URL:       w.URL,
			RateLimit: w.RateLimit,
			Payload:   msg.WebhookBody{Event: "sms_status", Data: p},
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

type PublishMOParams struct {
	AccountID   string       `json:"account_id"`
	SMSID       string       `json:"sms_id"`
	Message     string       `json:"message"`
	Recipient   string       `json:"recipient"`
	Sender      string       `json:"sender"`
	ReceivedAt  time.Time    `json:"received_at"`
	LastMessage *LastMessage `json:"last_message,omitempty"`
}

type LastMessage struct {
	Type        string   `json:"type"`
	ID          string   `json:"id"`
	Recipient   string   `json:"recipient"`
	Sender      string   `json:"sender"`
	Message     string   `json:"message"`
	MessageRef  string   `json:"message_ref"`
	Subject     string   `json:"subject,omitempty"`
	ContentURLS []string `json:"content_urls,omitempty"`
}

func (s *Webhook) PublishMO(p PublishMOParams, r *NoReply) error {
	webhooks, err := s.db.FindByEvent(p.AccountID, "sms_inbound")
	if err != nil {
		return err
	}

	for _, w := range webhooks {
		err = s.db.Publish(msg.WebhookMessageSpec{
			URL:       w.URL,
			RateLimit: w.RateLimit,
			Payload:   msg.WebhookBody{Event: "sms_inbound", Data: p},
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
