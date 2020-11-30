package client

import (
	rpc "github.com/burstsms/mtmo-tp/backend/sms/rpc/types"
)

type SendParams = rpc.SendParams
type SendReply = rpc.SendReply

func (c *Client) Send(params SendParams) (r *SendReply, err error) {
	r = &SendReply{}
	err = c.Call("Send", params, r)
	return r, err
}

type FindByIDParams = rpc.FindByIDParams
type FindByIDReply = rpc.FindByIDReply

func (c *Client) FindByID(params FindByIDParams) (r *FindByIDReply, err error) {
	r = &FindByIDReply{}
	err = c.Call("FindByID", params, r)
	return r, err
}

type MarkSentParams = rpc.MarkSentParams
type MarkFailedParams = rpc.MarkFailedParams

func (c *Client) MarkSent(params MarkSentParams) (err error) {
	err = c.Call("MarkSent", params, &NoReply{})
	return err
}

func (c *Client) MarkFailed(params MarkFailedParams) (err error) {
	err = c.Call("MarkFailed", params, &NoReply{})
	return err
}
