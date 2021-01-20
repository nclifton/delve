package queue

import (
	"context"

	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
)

type rabbitQueue struct {
	rabbit   rabbit.Conn
	postOpts rabbit.PublishOptions
}

type RabbitPublishOptions = rabbit.PublishOptions

func NewRabbitQueue(conn rabbit.Conn, postOpts rabbit.PublishOptions) Queue {
	return &rabbitQueue{rabbit: conn, postOpts: postOpts}
}

func (q *rabbitQueue) PostWebhook(ctx context.Context, msg PostWebhookMessage) error {
	opts := q.postOpts
	opts.Ctx = ctx
	return rabbit.Publish(q.rabbit, opts, msg)
}
