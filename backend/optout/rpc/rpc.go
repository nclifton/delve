package rpc

import (
	"context"
	"encoding/gob"

	"google.golang.org/grpc"

	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	mms "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/optout/rpc/types"
	sms "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
)

const Name = "OptOut"

type optOutDB interface {
	FindOptOutByLinkID(ctx context.Context, linkID string) (*types.OptOut, error)
	InsertOptOut(ctx context.Context, accountID, messageID, messageType, sender string) (*types.OptOut, error)
}

type webhookRPC interface {
	PublishOptOut(ctx context.Context, in *webhookpb.PublishOptOutParams, opts ...grpc.CallOption) (*webhookpb.NoReply, error)
}

type OptOutService struct {
	db           optOutDB
	webhookRPC   webhookRPC
	smsRPC       *sms.Client
	mmsRPC       *mms.Client
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

func NewService(postgresURL, optOutDomain string, webhook webhookpb.ServiceClient, sms *sms.Client, mms *mms.Client) (rpc.Service, error) {
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
			mmsRPC:       mms,
		},
	}

	return service, nil
}
