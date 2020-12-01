package client

import (
	"github.com/burstsms/mtmo-tp/backend/mms/rpc/types"
)

type SendParams types.SendParams
type SendReply types.SendReply

func (c *Client) Send(p SendParams) (r *SendReply, err error) {
	r = &SendReply{}
	err = c.Call("Send", p, r)
	return
}

type UpdateStatusParams = types.UpdateStatusParams

func (c *Client) UpdateStatus(p UpdateStatusParams) error {
	return c.Call("UpdateStatus", p, &types.NoReply{})
}

type FindByIDParams = types.FindByIDParams
type FindByIDReply = types.FindByIDReply

func (c *Client) FindByID(p FindByIDParams) (r *FindByIDReply, err error) {
	r = &FindByIDReply{}
	err = c.Call("FindByID", p, r)
	return
}
