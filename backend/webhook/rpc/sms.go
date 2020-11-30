package rpc

import (
	"github.com/burstsms/mtmo-tp/backend/webhook/db"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/types"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/msg"
)

func (s *Webhook) PublishSMSStatusUpdate(p types.PublishSMSStatusUpdateParams, r *types.NoReply) error {
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

func (s *Webhook) PublishMO(p types.PublishMOParams, r *types.NoReply) error {
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
