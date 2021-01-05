package queue

import (
	"context"

	"github.com/burstsms/mtmo-tp/backend/webhook/worker/msg"
)

type PostWebhookMessage = msg.WebhookMessageSpec

type Queue interface {
	PostWebhook(ctx context.Context, msg PostWebhookMessage) error
}
