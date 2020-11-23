package rpc

import (
	"encoding/gob"

	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	webhook "github.com/burstsms/mtmo-tp/backend/webhook/rpc/client"
)

const Name = "SMS"

type NoParams struct{}
type NoReply struct{}

type SMSService struct {
	db         *db
	webhookRPC *webhook.Client
	name       string
}

type Service struct {
	receiver *SMSService
}

func (s *Service) Name() string {
	return s.receiver.name
}

func (s *Service) Receiver() interface{} {
	return s.receiver
}

func NewService(postgresURL string, rabbitmq rabbit.Conn, webhook *webhook.Client) (rpc.Service, error) {
	gob.Register(map[string]interface{}{})
	db, err := NewDB(postgresURL, rabbitmq)
	if err != nil {
		return nil, err
	}
	service := &Service{
		receiver: &SMSService{db: db, name: Name, webhookRPC: webhook},
	}

	return service, nil
}
