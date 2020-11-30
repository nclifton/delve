package rpc

import (
	"encoding/gob"

	account "github.com/burstsms/mtmo-tp/backend/account/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	optOut "github.com/burstsms/mtmo-tp/backend/optout/rpc/client"
	tracklink "github.com/burstsms/mtmo-tp/backend/track_link/rpc/client"
	webhook "github.com/burstsms/mtmo-tp/backend/webhook/rpc/client"
)

const Name = "SMS"

type SMSService struct {
	db           *db
	webhookRPC   *webhook.Client
	accountRPC   *account.Client
	tracklinkRPC *tracklink.Client
	optOutRPC    *optOut.Client
	name         string
	features     SMSFeatures
}

type SMSFeatures struct {
	TrackLinkDomain  string
	OptOutLinkDomain string
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

func NewService(
	features SMSFeatures,
	postgresURL string,
	rabbitmq rabbit.Conn,
	webhook *webhook.Client,
	account *account.Client,
	tracklink *tracklink.Client,
	redisURL string,
	optOut *optOut.Client) (rpc.Service, error) {

	gob.Register(map[string]interface{}{})
	db, err := NewDB(postgresURL, rabbitmq, redisURL)
	if err != nil {
		return nil, err
	}
	service := &Service{
		receiver: &SMSService{
			db:           db,
			name:         Name,
			webhookRPC:   webhook,
			accountRPC:   account,
			tracklinkRPC: tracklink,
			features:     features,
			optOutRPC:    optOut,
		},
	}

	return service, nil
}
