package rpc

import (
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/redis"
)

type s3Svc interface {
	PutS3Content(content []byte, bucket, key string) error
}

// wrap the underlying connection so that exported methods
// can concistently drive database operations
type db struct {
	rabbit  rabbit.Conn
	opts    RabbitPublishOptions
	redis   *redis.Connection
	limiter *redis.Limiter
	s3      s3Svc
}

type RabbitPublishOptions = rabbit.PublishOptions

func (db *db) Publish(msg interface{}, routeKey string) error {
	publishOpts := RabbitPublishOptions{
		Exchange:     db.opts.Exchange,
		ExchangeType: db.opts.ExchangeType,
		RouteKey:     routeKey,
	}
	return rabbit.Publish(db.rabbit, publishOpts, msg)
}
