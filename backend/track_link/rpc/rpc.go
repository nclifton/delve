package rpc

import (
	"encoding/gob"

	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	mms "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
	sms "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
	webhook "github.com/burstsms/mtmo-tp/backend/webhook/rpc/client"
)

const Name = "TrackLink"

type TrackLinkService struct {
	db         *db
	name       string
	trackHost  string
	mmsRPC     *mms.Client
	smsRPC     *sms.Client
	webhookRPC *webhook.Client
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

// A silly comment
func NewService(postgresURL, trackHost string, mms *mms.Client, sms *sms.Client, webhook *webhook.Client) (rpc.Service, error) {
	gob.Register(map[string]interface{}{})
	db, err := NewDB(postgresURL)
	if err != nil {
		return nil, err
	}
	service := &Service{
		receiver: &TrackLinkService{
			db:         db,
			name:       Name,
			trackHost:  trackHost,
			mmsRPC:     mms,
			smsRPC:     sms,
			webhookRPC: webhook},
	}

	return service, nil
}
