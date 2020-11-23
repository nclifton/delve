package client

import (
	rpc "github.com/burstsms/mtmo-tp/backend/sms/rpc"
)

type QueueMOParams = rpc.QueueMOParams

func (c *Client) QueueMO(params QueueMOParams) (err error) {
	err = c.Call("QueueMO", params, &NoReply{})
	return err
}
