package rpc

import (
	"github.com/burstsms/mtmo-tp/backend/webhook/db"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/types"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/msg"
)

const eventMMSStatus = "mms_status"

func (s *Webhook) PublishMMSStatusUpdate(p types.PublishMMSStatusUpdateParams, r *types.NoReply) error {
	webhooks, err := s.db.FindByEvent(p.AccountID, eventMMSStatus)
	if err != nil {
		return err
	}

	for _, w := range webhooks {
		err = s.db.Publish(msg.WebhookMessageSpec{
			URL:       w.URL,
			RateLimit: w.RateLimit,
			Payload:   msg.WebhookBody{Event: eventMMSStatus, Data: p},
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
