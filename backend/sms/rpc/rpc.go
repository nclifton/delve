package rpc

import (
	"encoding/gob"

	account "github.com/burstsms/mtmo-tp/backend/account/rpc/client"
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
	accountRPC *account.Client
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

func NewService(postgresURL string, rabbitmq rabbit.Conn, webhook *webhook.Client, account *account.Client, redisURL string) (rpc.Service, error) {
	gob.Register(map[string]interface{}{})
	db, err := NewDB(postgresURL, rabbitmq, redisURL)
	if err != nil {
		return nil, err
	}
	service := &Service{
		receiver: &SMSService{db: db, name: Name, webhookRPC: webhook, accountRPC: account},
	}

	return service, nil
}
