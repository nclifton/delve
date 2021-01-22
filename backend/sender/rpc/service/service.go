package service

import (
	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/db"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
)

type senderImpl struct {
	log *logger.StandardLogger
	db  db.DB
	senderpb.UnimplementedServiceServer
}

func NewSenderService(db db.DB) senderpb.ServiceServer {
	return &senderImpl{
		log: logger.NewLogger(),
		db:  db,
	}
}
