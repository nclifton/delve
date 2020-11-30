package client

import (
	"github.com/burstsms/mtmo-tp/backend/optout/rpc/types"
)

type FindByLinkIDParams = types.FindByLinkIDParams
type FindByLinkIDReply = types.FindByLinkIDReply

func (c *Client) FindByLinkID(params FindByLinkIDParams) (r *FindByLinkIDReply, err error) {
	r = &FindByLinkIDReply{}
	err = c.Call(types.FindByLinkID, params, r)
	return r, err
}

type OptOutViaLinkParams = types.OptOutViaLinkParams
type OptOutViaLinkReply = types.OptOutViaLinkReply

func (c *Client) OptOutViaLink(params OptOutViaLinkParams) (r *OptOutViaLinkReply, err error) {
	r = &OptOutViaLinkReply{}
	err = c.Call(types.OptOutViaLink, params, r)
	return r, err
}

type GenerateOptOutLinkParams = types.GenerateOptOutLinkParams
type GenerateOptOutLinkReply = types.GenerateOptOutLinkReply

func (c *Client) GenerateOptOutLink(params GenerateOptOutLinkParams) (r *GenerateOptOutLinkReply, err error) {
	r = &GenerateOptOutLinkReply{}
	err = c.Call(types.GenerateOptOutLink, params, r)
	return r, err
}

type OptOutViaMsgParams = types.OptOutViaMsgParams

func (c *Client) OptOutViaMsg(params types.OptOutViaMsgParams) error {
	return c.Call(types.OptOutViaMsg, params, &types.NoReply{})
}
