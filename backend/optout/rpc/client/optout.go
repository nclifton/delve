package client

import (
	"github.com/burstsms/mtmo-tp/backend/optout/rpc"
	"github.com/burstsms/mtmo-tp/backend/optout/rpc/types"
)

type FindByLinkIDParams = types.FindByLinkIDParams
type FindByLinkIDReply = types.FindByLinkIDReply

func (c *Client) FindByLinkID(params FindByLinkIDParams) (r *FindByLinkIDReply, err error) {
	r = &FindByLinkIDReply{}
	err = c.Call("FindByLinkID", params, r)
	return r, err
}

type OptOutViaLinkParams = types.OptOutViaLinkParams
type OptOutViaLinkReply = types.OptOutViaLinkReply

func (c *Client) OptOutViaLink(params OptOutViaLinkParams) (r *OptOutViaLinkReply, err error) {
	r = &OptOutViaLinkReply{}
	err = c.Call("OptOutViaLink", params, r)
	return r, err
}

type GenerateOptOutLinkParams = rpc.GenerateOptOutLinkParams
type GenerateOptOutLinkReply = rpc.GenerateOptOutLinkReply

func (c *Client) GenerateOptOutLink(params GenerateOptOutLinkParams) (r *GenerateOptOutLinkReply, err error) {
	r = &GenerateOptOutLinkReply{}
	err = c.Call("GenerateOptOutLink", params, r)
	return r, err
}
