package rpc

import (
	"encoding/gob"

	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
)

const Name = "MMS"

type NoParams struct{}
type NoReply struct{}

type MMSService struct {
	db   *db
	name string
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

func NewService(postgresURL string) (rpc.Service, error) {
	gob.Register(map[string]interface{}{})
	db, err := NewDB(postgresURL)
	if err != nil {
		return nil, err
	}
	service := &Service{
		receiver: &MMSService{db: db, name: Name},
	}

	return service, nil
}
