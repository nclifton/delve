package client

import (
	"github.com/burstsms/mtmo-tp/backend/sms/rpc/types"
	rpc "github.com/burstsms/mtmo-tp/backend/sms/rpc/types"
)

type SMS = types.SMS

type SendParams = rpc.SendParams
type SendReply = rpc.SendReply

func (c *Client) Send(p SendParams) (r *SendReply, err error) {
	r = &SendReply{}
	err = c.Call("Send", p, r)
	return r, err
}

type FindByIDParams = rpc.FindByIDParams
type FindByIDReply = rpc.FindByIDReply

func (c *Client) FindByID(p FindByIDParams) (r *FindByIDReply, err error) {
	r = &FindByIDReply{}
	err = c.Call("FindByID", p, r)
	return r, err
}

type MarkSentParams = rpc.MarkSentParams
type MarkFailedParams = rpc.MarkFailedParams

func (c *Client) MarkSent(p MarkSentParams) error {
	return c.Call("MarkSent", p, &types.NoReply{})
}

func (c *Client) MarkFailed(p MarkFailedParams) error {
	return c.Call("MarkFailed", p, &types.NoReply{})
}
