package client

import (
	"github.com/burstsms/mtmo-tp/backend/sms/rpc/types"
	rpc "github.com/burstsms/mtmo-tp/backend/sms/rpc/types"
)

type QueueMOParams = rpc.QueueMOParams

func (c *Client) QueueMO(params QueueMOParams) error {
	return c.Call("QueueMO", params, &types.NoReply{})
}

type ProcessMOParams = rpc.ProcessMOParams

func (c *Client) ProcessMO(params ProcessMOParams) error {
	return c.Call("ProcessMO", params, &types.NoReply{})
}
