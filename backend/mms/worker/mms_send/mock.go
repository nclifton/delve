package mms_send

import (
	mm7RPC "github.com/burstsms/mtmo-tp/backend/mm7/rpc"
	mmsRPC "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
)

type mockMM7RPCClient struct {
	err error
}

func (m mockMM7RPCClient) Send(p mm7RPC.MM7SendParams) (r *mm7RPC.NoReply, err error) {
	return nil, m.err
}

type mockMMSRPCClient struct {
	err error
}

func (m mockMMSRPCClient) UpdateStatus(p mmsRPC.UpdateStatusParams) error {
	return m.err
}
