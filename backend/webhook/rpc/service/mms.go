package service

import (
	"context"
	"log"
	"time"

	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/queue"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/msg"
)

func (s *webhookImpl) PublishMMSStatusUpdate(ctx context.Context, p *webhookpb.PublishMMSStatusUpdateParams) (*webhookpb.NoReply, error) {
	webhooks, err := s.db.FindWebhookByEvent(ctx, p.AccountId, EventMMSStatus)
	if err != nil {
		log.Printf("webhook for account id %s with event %s was not found", p.AccountId, EventMMSStatus)
		return &webhookpb.NoReply{}, err
	}

	for _, w := range webhooks {
		err = s.queue.PostWebhook(ctx, queue.PostWebhookMessage{
			URL:       w.URL,
			RateLimit: int(w.RateLimit),
			Payload: msg.WebhookBody{
				Event: EventMMSStatus,
				Data: PublishStatusData{
					MMS_id:            p.MMSId,
					Message_ref:       p.MessageRef,
					Recipient:         p.Recipient,
					Sender:            p.Sender,
					Status:            p.Status,
					Status_updated_at: p.StatusUpdatedAt.AsTime().Format(time.RFC3339),
				}},
		})
		if err != nil {
			return &webhookpb.NoReply{}, err
		}
	}

	return &webhookpb.NoReply{}, nil
}
