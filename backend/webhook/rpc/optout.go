package rpc

import (
	"github.com/burstsms/mtmo-tp/backend/webhook/db"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/types"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/msg"
)

func (s *Webhook) PublishOptOut(p types.PublishOptOutParams, r *types.NoReply) error {
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
