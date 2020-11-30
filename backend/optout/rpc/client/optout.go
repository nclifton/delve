package client

import (
	rpc "github.com/burstsms/mtmo-tp/backend/optout/rpc"
)

type FindByLinkIDParams = rpc.FindByLinkIDParams
type FindByLinkIDReply = rpc.FindByLinkIDReply

func (c *Client) FindByLinkID(params FindByLinkIDParams) (r *FindByLinkIDReply, err error) {
	r = &FindByLinkIDReply{}
	err = c.Call("FindByLinkID", params, r)
	return r, err
}

type OptOutViaLinkParams = rpc.OptOutViaLinkParams
type OptOutViaLinkReply = rpc.OptOutViaLinkReply

func (c *Client) OptOutViaLink(params OptOutViaLinkParams) (r *OptOutViaLinkReply, err error) {
	r = &OptOutViaLinkReply{}
	err = c.Call("OptOutViaLink", params, r)
	return r, err
}

type GenerateOptoutLinkParams = rpc.GenerateOptoutLinkParams
type GenerateOptoutLinkReply = rpc.GenerateOptoutLinkReply

func (c *Client) GenerateOptoutLink(params GenerateOptoutLinkParams) (r *GenerateOptoutLinkReply, err error) {
	r = &GenerateOptoutLinkReply{}
	err = c.Call("GenerateOptoutLink", params, r)
	return r, err
}
