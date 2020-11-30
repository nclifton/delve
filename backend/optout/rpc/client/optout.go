package client

import (
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

type GenerateOptoutLinkParams = types.GenerateOptoutLinkParams
type GenerateOptoutLinkReply = types.GenerateOptoutLinkReply

func (c *Client) GenerateOptoutLink(params GenerateOptoutLinkParams) (r *GenerateOptoutLinkReply, err error) {
	r = &GenerateOptoutLinkReply{}
	err = c.Call("GenerateOptoutLink", params, r)
	return r, err
}
