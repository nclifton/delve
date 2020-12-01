package client

import (
	"github.com/burstsms/mtmo-tp/backend/mm7/rpc/types"
)

type SendParams = types.SendParams

func (c *Client) Send(p SendParams) error {
	return c.Call("Send", p, &types.NoReply{})
}

type ProviderSpecParams = types.ProviderSpecParams
type ProviderSpecReply = types.ProviderSpecReply

func (c *Client) ProviderSpec(p ProviderSpecParams) (r *ProviderSpecReply, err error) {
	r = &ProviderSpecReply{}
	err = c.Call("ProviderSpec", p, r)
	return
}

type UpdateStatusParams = types.UpdateStatusParams

func (c *Client) UpdateStatus(p UpdateStatusParams) error {
	return c.Call("UpdateStatus", p, &types.NoReply{})
}

type DLRParams = types.DLRParams

func (c *Client) DLR(p DLRParams) error {
	return c.Call("DLR", p, &types.NoReply{})
}

type DeliverParams = types.DeliverParams

func (c *Client) Deliver(p DeliverParams) error {
	return c.Call("Deliver", p, &types.NoReply{})
}

type GetCachedContentParams = types.GetCachedContentParams
type GetCachedContentReply = types.GetCachedContentReply

func (c *Client) GetCachedContent(p GetCachedContentParams) (r *GetCachedContentReply, err error) {
	r = &GetCachedContentReply{}
	err = c.Call("GetCachedContent", p, r)
	return
}

type CheckRateLimitParams = types.CheckRateLimitParams
type CheckRateLimitReply = types.CheckRateLimitReply

func (c *Client) CheckRateLimit(p CheckRateLimitParams) (r *CheckRateLimitReply, err error) {
	r = &CheckRateLimitReply{}
	err = c.Call("CheckRateLimit", p, r)
	return
}
