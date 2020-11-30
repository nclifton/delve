package rpc

import (
	"context"
	types "github.com/burstsms/mtmo-tp/backend/optout/rpc/types"
)

type mockDB struct {
	optOut types.OptOut
	err    error
}

func (m mockDB) FindOptOutByLinkID(ctx context.Context, linkID string) (*types.OptOut, error) {
	return &m.optOut, m.err
}

func (m mockDB) InsertOptOut(ctx context.Context, accountID, messageID, messageType string) (*types.OptOut, error) {
	return &m.optOut, m.err
}
