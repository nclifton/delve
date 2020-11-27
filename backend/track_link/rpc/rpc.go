package rpc

import (
	"encoding/gob"

	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
)

const Name = "track_link"

type NoParams struct{}
type NoReply struct{}

type TrackLinkService struct {
	db        *db
	name      string
	trackHost string
}

type Service struct {
	receiver *TrackLinkService
}

func (s *Service) Name() string {
	return s.receiver.name
}

func (s *Service) Receiver() interface{} {
	return s.receiver
}

func NewService(postgresURL, trackHost string) (rpc.Service, error) {
	gob.Register(map[string]interface{}{})
	db, err := NewDB(postgresURL)
	if err != nil {
		return nil, err
	}
	service := &Service{
		receiver: &TrackLinkService{db: db, name: Name, trackHost: trackHost},
	}

	return service, nil
}
