package service

import (
	"context"
	"time"

	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/post/msg"
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
			Payload: msg.WebhookBody{
				Event: EventLinkHitStatus,
				Data: PublishLinkHitData{
					URL:       p.URL,
					Hits:      int(p.Hits),
					Timestamp: p.Timestamp.AsTime().Format(time.RFC3339),
					Source_message: PublishMessageData{
						Type:         p.SourceMessage.Type,
						Id:           p.SourceMessage.Id,
						Recipient:    p.SourceMessage.Recipient,
						Sender:       p.SourceMessage.Sender,
						Subject:      p.SourceMessage.Subject,
						Message:      p.SourceMessage.Message,
						Content_urls: p.SourceMessage.ContentURLs,
						Message_ref:  p.SourceMessage.MessageRef,
					}}}})
		if err != nil {
			return &webhookpb.NoReply{}, err
		}
	}

	return &webhookpb.NoReply{}, nil
}
