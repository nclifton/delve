package fakemm7dldrworker

import (
	mm7RPC "github.com/burstsms/mtmo-tp/backend/mm7/rpc/client"
)

type MockRPCClient struct {
	MediaStoreReply mm7RPC.MediaStoreReply
	DLRReply        mm7RPC.NoReply
	DeliverReply    mm7RPC.NoReply
	Error           error
}

func (m MockRPCClient) Store(params mm7RPC.MediaStoreParams) (r *mm7RPC.MediaStoreReply, err error) {
	return &m.MediaStoreReply, m.Error
}

func (m MockRPCClient) DLR(params mm7RPC.DLRParams) (r *mm7RPC.NoReply, err error) {
	return nil, m.Error
}

func (m MockRPCClient) Deliver(params mm7RPC.DeliverParams) (r *mm7RPC.NoReply, err error) {
	return nil, m.Error
}
