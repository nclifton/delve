package rpc

import (
	"context"
	"encoding/gob"

	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	"github.com/burstsms/mtmo-tp/backend/optout/rpc/types"
	sms "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
	webhook "github.com/burstsms/mtmo-tp/backend/webhook/rpc/client"
)

const Name = "OptOut"

type optOutDB interface {
	FindOptOutByLinkID(ctx context.Context, linkID string) (*types.OptOut, error)
	InsertOptOut(ctx context.Context, accountID, messageID, messageType, sender string) (*types.OptOut, error)
}

type webhookRPC interface {
	PublishOptOut(p webhook.PublishOptOutParams) error
}

type OptOutService struct {
	db           optOutDB
	webhookRPC   webhookRPC
	smsRPC       *sms.Client
	name         string
	optOutDomain string
}

type Service struct {
	receiver *OptOutService
}

func (s *Service) Name() string {
	return s.receiver.name
}

func (s *Service) Receiver() interface{} {
	return s.receiver
}

func NewService(postgresURL, optOutDomain string, webhook *webhook.Client, sms *sms.Client) (rpc.Service, error) {
	gob.Register(map[string]interface{}{})
	db, err := NewDB(postgresURL)
	if err != nil {
		return nil, err
	}
	service := &Service{
		receiver: &OptOutService{
			db:           db,
			name:         Name,
			optOutDomain: optOutDomain,
			webhookRPC:   webhook,
			smsRPC:       sms,
		},
	}

	return service, nil
}
