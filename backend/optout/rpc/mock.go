package rpc

import "context"

type mockDB struct {
	optOut OptOut
	err    error
}

func (m mockDB) FindOptOutByLinkID(ctx context.Context, linkID string) (*OptOut, error) {
	return &m.optOut, m.err
}

func (m mockDB) InsertOptOut(ctx context.Context, accountID, messageID, messageType string) (*OptOut, error) {
	return &m.optOut, m.err
}
