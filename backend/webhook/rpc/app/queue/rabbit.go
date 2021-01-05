package queue

import (
	"context"

	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
)

type rabbitQueue struct {
	rabbit rabbit.Conn
	opts   rabbit.PublishOptions
}

type RabbitPublishOptions = rabbit.PublishOptions

func NewRabbitQueue(conn rabbit.Conn, opts rabbit.PublishOptions) Queue {
	return &rabbitQueue{rabbit: conn, opts: opts}
}

func (q *rabbitQueue) PostWebhook(ctx context.Context, msg PostWebhookMessage) error {
	opts := q.opts
	opts.Ctx = ctx
	return rabbit.Publish(q.rabbit, opts, msg)
}
