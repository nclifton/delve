package client

import (
	rpc "github.com/burstsms/mtmo-tp/backend/sms/rpc"
)

type QueueDLRParams = rpc.QueueDLRParams

func (c *Client) QueueDLR(params QueueDLRParams) (err error) {
	err = c.Call("QueueDLR", params, &NoReply{})
	return err
}
