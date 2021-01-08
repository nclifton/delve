package service

import (
	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/app/db"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/app/queue"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
)

type webhookImpl struct {
	log   *logger.StandardLogger
	db    db.DB
	queue queue.Queue
	webhookpb.UnimplementedServiceServer
}

func NewWebhookService(db db.DB, queue queue.Queue) webhookpb.ServiceServer {
	return &webhookImpl{
		log:   logger.NewLogger(),
		db:    db,
		queue: queue,
	}
}
