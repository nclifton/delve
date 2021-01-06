package service

import (
	"context"

	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/msg"
)

func (s *webhookImpl) PublishOptOut(ctx context.Context, p *webhookpb.PublishOptOutParams) (*webhookpb.NoReply, error) {
	webhooks, err := s.db.FindWebhookByEvent(ctx, p.AccountId, EventOptOutStatus)
	if err != nil {
		s.log.Error(ctx, "FindWebhookByEvent", err.Error())
		return &webhookpb.NoReply{}, err
	}

	for _, w := range webhooks {
		err = s.queue.PostWebhook(ctx, msg.WebhookMessageSpec{
			URL:       w.URL,
			RateLimit: int(w.RateLimit),
			Payload:   msg.WebhookBody{Event: EventOptOutStatus, Data: p},
		})
		if err != nil {
			s.log.Error(ctx, "PostWebhook", err.Error())
			return &webhookpb.NoReply{}, err
		}
	}

	return &webhookpb.NoReply{}, nil
}
