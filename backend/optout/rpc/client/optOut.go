package client

import (
	"github.com/burstsms/mtmo-tp/backend/optout/rpc/types"
)

type FindByLinkIDParams = types.FindByLinkIDParams
type FindByLinkIDReply = types.FindByLinkIDReply

func (c *Client) FindByLinkID(p FindByLinkIDParams) (r *FindByLinkIDReply, err error) {
	r = &FindByLinkIDReply{}
	err = c.Call("FindByLinkID", p, r)
	return r, err
}

type OptOutViaLinkParams = types.OptOutViaLinkParams
type OptOutViaLinkReply = types.OptOutViaLinkReply

func (c *Client) OptOutViaLink(p OptOutViaLinkParams) (r *OptOutViaLinkReply, err error) {
	r = &OptOutViaLinkReply{}
	err = c.Call("OptOutViaLink", p, r)
	return r, err
}

type GenerateOptOutLinkParams = types.GenerateOptOutLinkParams
type GenerateOptOutLinkReply = types.GenerateOptOutLinkReply

func (c *Client) GenerateOptOutLink(p GenerateOptOutLinkParams) (r *GenerateOptOutLinkReply, err error) {
	r = &GenerateOptOutLinkReply{}
	err = c.Call("GenerateOptOutLink", p, r)
	return r, err
}

type OptOutViaMsgParams = types.OptOutViaMsgParams

func (c *Client) OptOutViaMsg(p types.OptOutViaMsgParams) error {
	return c.Call("OptOutViaMsg", p, &types.NoReply{})
}
