package service

import (
	"github.com/burstsms/mtmo-tp/backend/account/rpc/accountpb"
	"github.com/burstsms/mtmo-tp/backend/account/rpc/db"
	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/redis"
)

type accountImpl struct {
	log *logger.StandardLogger
	accountpb.UnimplementedServiceServer
	db    db.DB
	redis *redis.Connection
}

func NewAccountService(db db.DB, redis *redis.Connection) accountpb.ServiceServer {
	return &accountImpl{
		log:   logger.NewLogger(),
		db:    db,
		redis: redis,
	}
}
