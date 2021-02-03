package rpc

import (
	"context"

	"google.golang.org/grpc"

	mms "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
	types "github.com/burstsms/mtmo-tp/backend/optout/rpc/types"
	sms "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
)

type mockDB struct {
	optOut types.OptOut
	err    error
}

func (m mockDB) FindOptOutByLinkID(ctx context.Context, linkID string) (*types.OptOut, error) {
	return &m.optOut, m.err
}

func (m mockDB) InsertOptOut(ctx context.Context, accountID, messageID, messageType, sender string) (*types.OptOut, error) {
	return &m.optOut, m.err
}

type mockWebhookRPC struct {
	reply *webhookpb.NoReply
	err   error
}

func (m mockWebhookRPC) PublishOptOut(ctx context.Context, p *webhookpb.PublishOptOutParams, opts ...grpc.CallOption) (*webhookpb.NoReply, error) {
	return m.reply, m.err
}

type mockSMSRPC struct {
	findByIDReply sms.FindByIDReply
	err           error
}

func (m mockSMSRPC) FindByID(p sms.FindByIDParams) (r *sms.FindByIDReply, err error) {
	return &m.findByIDReply, err
}

type mockMMSRPC struct {
	findByIDReply mms.FindByIDReply
	err           error
}

func (m mockMMSRPC) FindByID(p mms.FindByIDParams) (r *mms.FindByIDReply, err error) {
	return &m.findByIDReply, err
}
