package service

import (
	"context"

	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/msg"
)

func (s *webhookImpl) PublishLinkHit(ctx context.Context, p *webhookpb.PublishLinkHitParams) (*webhookpb.NoReply, error) {
	webhooks, err := s.db.FindWebhookByEvent(ctx, p.AccountId, EventLinkHitStatus)
	if err != nil {
		return &webhookpb.NoReply{}, err
	}

	for _, w := range webhooks {
		err = s.queue.PostWebhook(ctx, msg.WebhookMessageSpec{
			URL:       w.URL,
			RateLimit: int(w.RateLimit),
			Payload:   msg.WebhookBody{Event: EventLinkHitStatus, Data: p},
		})
		if err != nil {
			return &webhookpb.NoReply{}, err
		}
	}

	return &webhookpb.NoReply{}, nil
}
