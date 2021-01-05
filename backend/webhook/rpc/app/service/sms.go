package service

import (
	"context"

	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/msg"
)

func (s *webhookImpl) PublishSMSStatusUpdate(ctx context.Context, p *webhookpb.PublishSMSStatusUpdateParams) (*webhookpb.NoReply, error) {
	webhooks, err := s.db.FindWebhookByEvent(ctx, p.AccountId, EventSMSStatus)
	if err != nil {
		return &webhookpb.NoReply{}, err
	}

	for _, w := range webhooks {
		err = s.queue.PostWebhook(ctx, msg.WebhookMessageSpec{
			URL:       w.URL,
			RateLimit: int(w.RateLimit),
			Payload:   msg.WebhookBody{Event: EventSMSStatus, Data: p},
		})
		if err != nil {
			return &webhookpb.NoReply{}, err
		}
	}

	return &webhookpb.NoReply{}, nil
}

func (s *webhookImpl) PublishMO(ctx context.Context, p *webhookpb.PublishMOParams) (*webhookpb.NoReply, error) {
	webhooks, err := s.db.FindWebhookByEvent(ctx, p.AccountId, EventMOStatus)
	if err != nil {
		return &webhookpb.NoReply{}, err
	}

	for _, w := range webhooks {
		err = s.queue.PostWebhook(ctx, msg.WebhookMessageSpec{
			URL:       w.URL,
			RateLimit: int(w.RateLimit),
			Payload:   msg.WebhookBody{Event: EventMOStatus, Data: p},
		})
		if err != nil {
			return &webhookpb.NoReply{}, err
		}
	}

	return &webhookpb.NoReply{}, nil
}
