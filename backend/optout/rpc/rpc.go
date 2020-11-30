package rpc

import (
	"context"
	"encoding/gob"

	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	webhook "github.com/burstsms/mtmo-tp/backend/webhook/rpc/client"
)

const Name = "OptOut"

type NoParams struct{}
type NoReply struct{}

type optOutDB interface {
	FindOptOutByLinkID(ctx context.Context, linkID string) (*OptOut, error)
	InsertOptOut(ctx context.Context, accountID, messageID, messageType string) (*OptOut, error)
}

type OptOutService struct {
	db         optOutDB
	webhookRPC *webhook.Client
	name       string
	trackHost  string
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

func NewService(postgresURL, trackHost string, webhook *webhook.Client) (rpc.Service, error) {
	gob.Register(map[string]interface{}{})
	db, err := NewDB(postgresURL)
	if err != nil {
		return nil, err
	}
	service := &Service{
		receiver: &OptOutService{db: db, name: Name, trackHost: trackHost, webhookRPC: webhook},
	}

	return service, nil
}
