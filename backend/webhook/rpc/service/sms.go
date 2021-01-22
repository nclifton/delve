package service

import (
	"context"
	"time"

	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/post/msg"
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
			Payload: msg.WebhookBody{
				Event: EventSMSStatus,
				Data: PublishStatusData{
					SMS_id:            p.SMSId,
					Message_ref:       p.MessageRef,
					Recipient:         p.Recipient,
					Sender:            p.Sender,
					Status:            p.Status,
					Status_updated_at: p.StatusUpdatedAt.AsTime().Format(time.RFC3339),
				}}})
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
			Payload: msg.WebhookBody{
				Event: EventMOStatus,
				Data: PublishMOData{
					SMS_id:    p.SMSId,
					Recipient: p.Recipient,
					Sender:    p.Sender,
					Message:   p.Message,
					Timestamp: p.ReceivedAt.AsTime().Format(time.RFC3339),
					Last_message: PublishMessageData{
						Type:         p.LastMessage.Type,
						Id:           p.LastMessage.Id,
						Recipient:    p.LastMessage.Recipient,
						Sender:       p.LastMessage.Sender,
						Subject:      p.LastMessage.Subject,
						Message:      p.LastMessage.Message,
						Content_urls: p.LastMessage.ContentURLs,
						Message_ref:  p.LastMessage.MessageRef,
					},
				}}})
		if err != nil {
			return &webhookpb.NoReply{}, err
		}
	}

	return &webhookpb.NoReply{}, nil
}
