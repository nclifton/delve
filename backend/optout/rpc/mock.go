package rpc

import (
	"context"

	types "github.com/burstsms/mtmo-tp/backend/optout/rpc/types"
	webhook "github.com/burstsms/mtmo-tp/backend/webhook/rpc/client"
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
	err error
}

func (m mockWebhookRPC) PublishOptOut(p webhook.PublishOptOutParams) error {
	return m.err
}
