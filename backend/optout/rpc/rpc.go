package rpc

import (
	"encoding/gob"

	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	webhook "github.com/burstsms/mtmo-tp/backend/webhook/rpc/client"
)

const Name = "OptOut"

type NoParams struct{}
type NoReply struct{}

type OptOutService struct {
	db         *db
	webhookRPC *webhook.Client
	name       string
}

type Service struct {
	receiver *OptOutService
}

func (s *Service) Name() string {
	return s.receiver.name
}

func (s *Service) Receiver() interface{} {
	return s.receiver
}

func NewService(postgresURL string, webhook *webhook.Client, redisURL string) (rpc.Service, error) {
	gob.Register(map[string]interface{}{})
	db, err := NewDB(postgresURL, redisURL)
	if err != nil {
		return nil, err
	}
	service := &Service{
		receiver: &OptOutService{db: db, name: Name, webhookRPC: webhook},
	}

	return service, nil
}
