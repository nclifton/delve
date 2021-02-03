package service

import (
	"context"
	"time"

	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/post/msg"
)

func (s *webhookImpl) PublishOptOut(ctx context.Context, p *webhookpb.PublishOptOutParams) (*webhookpb.NoReply, error) {
	webhooks, err := s.db.FindWebhookByEvent(ctx, p.AccountId, EventOptOutStatus)
	if err != nil {
		s.log.Error(ctx, "FindWebhookByEvent", err.Error())
		return &webhookpb.NoReply{}, err
	}

	for _, w := range webhooks {
		messageData := PublishMessageData{}
		if p.SourceMessage != nil {
			messageData = PublishMessageData{
				Type:        p.SourceMessage.Type,
				Id:          p.SourceMessage.Id,
				Recipient:   p.SourceMessage.Recipient,
				Sender:      p.SourceMessage.Sender,
				Message:     p.SourceMessage.Message,
				Message_ref: p.SourceMessage.MessageRef,
			}
		}

		if err := s.queue.PostWebhook(ctx, msg.WebhookMessageSpec{
			URL:       w.URL,
			RateLimit: int(w.RateLimit),
			Payload: msg.WebhookBody{
				Event: EventOptOutStatus,
				Data: PublishOptOutData{
					Source:         p.Source,
					Timestamp:      p.Timestamp.AsTime().Format(time.RFC3339),
					Source_message: messageData,
				}},
		}); err != nil {
			s.log.Error(ctx, "PostWebhook", err.Error())
			return &webhookpb.NoReply{}, err
		}
	}

	return &webhookpb.NoReply{}, nil
}
