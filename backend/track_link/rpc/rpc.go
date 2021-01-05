package rpc

import (
	"encoding/gob"

	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	mms "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
	sms "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
)

const Name = "TrackLink"

type TrackLinkService struct {
	db          *db
	name        string
	trackDomain string
	mmsRPC      *mms.Client
	smsRPC      *sms.Client
	webhookRPC  webhookpb.ServiceClient
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
func NewService(postgresURL, trackDomain string, mms *mms.Client, sms *sms.Client, webhook webhookpb.ServiceClient) (rpc.Service, error) {
	gob.Register(map[string]interface{}{})
	db, err := NewDB(postgresURL)
	if err != nil {
		return nil, err
	}
	service := &Service{
		receiver: &TrackLinkService{
			db:          db,
			name:        Name,
			trackDomain: trackDomain,
			mmsRPC:      mms,
			smsRPC:      sms,
			webhookRPC:  webhook},
	}

	return service, nil
}
