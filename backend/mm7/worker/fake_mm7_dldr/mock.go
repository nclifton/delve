package fakemm7dldrworker

import (
	mm7RPC "github.com/burstsms/mtmo-tp/backend/mm7/rpc/client"
)

type MockRPCClient struct {
	MediaStoreReply mm7RPC.MediaStoreReply
	Error           error
}

func (m MockRPCClient) Store(params mm7RPC.MediaStoreParams) (r *mm7RPC.MediaStoreReply, err error) {
	return &m.MediaStoreReply, m.Error
}

func (m MockRPCClient) DLR(params mm7RPC.DLRParams) error {
	return m.Error
}

func (m MockRPCClient) Deliver(params mm7RPC.DeliverParams) error {
	return m.Error
}
