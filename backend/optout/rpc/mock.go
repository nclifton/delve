package rpc

import (
	"context"

	"google.golang.org/grpc"

	types "github.com/burstsms/mtmo-tp/backend/optout/rpc/types"
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
