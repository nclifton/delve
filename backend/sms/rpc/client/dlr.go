package client

import (
	"github.com/burstsms/mtmo-tp/backend/sms/rpc/types"
	rpc "github.com/burstsms/mtmo-tp/backend/sms/rpc/types"
)

type QueueDLRParams = rpc.QueueDLRParams

func (c *Client) QueueDLR(params QueueDLRParams) error {
	return c.Call("QueueDLR", params, &types.NoReply{})
}

type ProcessDLRParams = rpc.ProcessDLRParams

func (c *Client) ProcessDLR(params ProcessDLRParams) error {
	return c.Call("ProcessDLR", params, &types.NoReply{})
}
