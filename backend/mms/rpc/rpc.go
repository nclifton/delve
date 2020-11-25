package rpc

import (
	"encoding/gob"

	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	webhook "github.com/burstsms/mtmo-tp/backend/webhook/rpc/client"
)

const Name = "MMS"

type webhookSvc interface {
	PublishMMSStatusUpdate(params webhook.PublishMMSStatusUpdateParams) error
}

type ConfigSvc struct {
	Webhook webhookSvc
}

type NoParams struct{}
type NoReply struct{}

type MMSService struct {
	db   *db
	name string
	svc  ConfigSvc
}

type Service struct {
	receiver *MMSService
}

func (s *Service) Name() string {
	return s.receiver.name
}

func (s *Service) Receiver() interface{} {
	return s.receiver
}

func NewService(postgresURL string, rabbitmq rabbit.Conn, opts rabbit.PublishOptions, svc ConfigSvc) (rpc.Service, error) {
	gob.Register(map[string]interface{}{})
	db, err := NewDB(postgresURL, rabbitmq, opts)
	if err != nil {
		return nil, err
	}
	service := &Service{
		receiver: &MMSService{db: db, name: Name, svc: svc},
	}

	return service, nil
}
