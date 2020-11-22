package rpc

import (
	"encoding/gob"

	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
)

const Name = "SMS"

type NoParams struct{}
type NoReply struct{}

type SMSService struct {
	db   *db
	name string
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

func NewService(postgresURL string, rabbitmq rabbit.Conn, opts RabbitPublishOptions) (rpc.Service, error) {
	gob.Register(map[string]interface{}{})
	db, err := NewDB(postgresURL, rabbitmq, opts)
	if err != nil {
		return nil, err
	}
	service := &Service{
		receiver: &SMSService{db: db, name: Name},
	}

	return service, nil
}
