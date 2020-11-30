package client

import (
	"github.com/burstsms/mtmo-tp/backend/mms/rpc"
)

type SendParams rpc.SendParams
type SendReply = rpc.SendReply

func (c *Client) Send(p SendParams) (r *SendReply, err error) {
	r = &SendReply{}
	err = c.Call("Send", p, r)
	return
}

type UpdateStatusParams rpc.UpdateStatusParams

func (c *Client) UpdateStatus(p UpdateStatusParams) (err error) {
	return c.Call("UpdateStatus", p, &NoReply{})
}
