package client

import (
	"github.com/burstsms/mtmo-tp/backend/mm7/rpc/types"
)

type SendParams = types.SendParams

func (c *Client) Send(p SendParams) (r *NoReply, err error) {
	r = &NoReply{}
	err = c.Call("Send", SendParams{
		ID:          p.ID,
		Subject:     p.Subject,
		Message:     p.Message,
		Sender:      p.Sender,
		Recipient:   p.Recipient,
		ContentURLs: p.ContentURLs,
		ProviderKey: p.ProviderKey,
	}, r)
	return
}

type ProviderSpecParams = types.ProviderSpecParams
type ProviderSpecReply = types.ProviderSpecReply

func (c *Client) ProviderSpec(p ProviderSpecParams) (r *ProviderSpecReply, err error) {
	r = &ProviderSpecReply{}
	err = c.Call("ProviderSpec", ProviderSpecParams{
		ProviderKey: p.ProviderKey,
	}, r)
	return
}

type UpdateStatusParams = types.UpdateStatusParams

func (c *Client) UpdateStatus(p UpdateStatusParams) error {
	r := &NoReply{}

	return c.Call("UpdateStatus", UpdateStatusParams{
		ID:          p.ID,
		MessageID:   p.MessageID,
		Status:      p.Status,
		Description: p.Description,
	}, r)
}

type DLRParams = types.DLRParams

func (c *Client) DLR(p DLRParams) (r *NoReply, err error) {
	r = &NoReply{}
	err = c.Call("DLR", DLRParams{
		ID:          p.ID,
		Status:      p.Status,
		Description: p.Description,
	}, r)
	return
}

type DeliverParams = types.DeliverParams

func (c *Client) Deliver(p DeliverParams) (r *NoReply, err error) {
	r = &NoReply{}
	err = c.Call("Deliver", DeliverParams{
		Subject:     p.Subject,
		Message:     p.Message,
		Sender:      p.Sender,
		Recipient:   p.Recipient,
		ContentURLs: p.ContentURLs,
		ProviderKey: p.ProviderKey,
	}, r)
	return
}

type GetCachedContentParams = types.GetCachedContentParams
type GetCachedContentReply = types.GetCachedContentReply

func (c *Client) GetCachedContent(p GetCachedContentParams) (r *GetCachedContentReply, err error) {
	r = &GetCachedContentReply{}
	err = c.Call("GetCachedContent", GetCachedContentParams{
		ContentURL: p.ContentURL,
	}, r)
	return
}

type CheckRateLimitParams = types.CheckRateLimitParams
type CheckRateLimitReply = types.CheckRateLimitReply

func (c *Client) CheckRateLimit(p CheckRateLimitParams) (r *CheckRateLimitReply, err error) {
	r = &CheckRateLimitReply{}
	err = c.Call("CheckRateLimit", CheckRateLimitParams{
		ProviderKey: p.ProviderKey,
	}, r)
	return
}
