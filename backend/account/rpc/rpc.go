package rpc

import (
	"encoding/gob"

	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
)

const Name = "Account"

type NoParams struct{}
type NoReply struct{}

type AccountService struct {
	db   *db
	name string
}

type Service struct {
	receiver *AccountService
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
		receiver: &AccountService{db: db, name: Name},
	}

	return service, nil
}
