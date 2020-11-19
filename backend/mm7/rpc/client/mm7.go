package client

import (
	mrpc "github.com/burstsms/mtmo-tp/backend/mm7/rpc"
)

type PingResponse = mrpc.PingResponse

func (c *Client) Ping() (r *mrpc.PingResponse, err error) {
	r = &PingResponse{}
	err = c.Call("Ping", mrpc.NoParams{}, r)
	return
}

type MM7SendParams = mrpc.MM7SendParams

func (c *Client) Send(p MM7SendParams) (r *NoReply, err error) {
	r = &NoReply{}
	err = c.Call("Send", MM7SendParams{
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

type MM7ProviderSpecParams = mrpc.MM7ProviderSpecParams
type MM7ProviderSpecReply = mrpc.MM7ProviderSpecReply

func (c *Client) ProviderSpec(p MM7ProviderSpecParams) (r *MM7ProviderSpecReply, err error) {
	r = &MM7ProviderSpecReply{}
	err = c.Call("ProviderSpec", MM7ProviderSpecParams{
		ProviderKey: p.ProviderKey,
	}, r)
	return
}

type MM7UpdateStatusParams = mrpc.MM7UpdateStatusParams

func (c *Client) UpdateStatus(p MM7UpdateStatusParams) (r *NoReply, err error) {
	r = &NoReply{}
	err = c.Call("UpdateStatus", MM7UpdateStatusParams{
		ID:          p.ID,
		Status:      p.Status,
		Description: p.Description,
	}, r)
	return
}

type MM7DLRParams = mrpc.MM7DLRParams

func (c *Client) DLR(p MM7DLRParams) (r *NoReply, err error) {
	r = &NoReply{}
	err = c.Call("DLR", MM7DLRParams{
		ID:          p.ID,
		Status:      p.Status,
		Description: p.Description,
	}, r)
	return
}

type MM7DeliverParams = mrpc.MM7DeliverParams

func (c *Client) Deliver(p MM7DeliverParams) (r *NoReply, err error) {
	r = &NoReply{}
	err = c.Call("Deliver", MM7DeliverParams{
		Subject:     p.Subject,
		Message:     p.Message,
		Sender:      p.Sender,
		Recipient:   p.Recipient,
		ContentURLs: p.ContentURLs,
		ProviderKey: p.ProviderKey,
	}, r)
	return
}

type MM7GetCachedContentParams = mrpc.MM7GetCachedContentParams
type MM7GetCachedContentReply = mrpc.MM7GetCachedContentReply

func (c *Client) GetCachedContent(p MM7GetCachedContentParams) (r *MM7GetCachedContentReply, err error) {
	r = &MM7GetCachedContentReply{}
	err = c.Call("GetCachedContent", MM7GetCachedContentParams{
		ContentURL: p.ContentURL,
	}, r)
	return
}

type MM7CheckRateLimitParams = mrpc.MM7CheckRateLimitParams
type MM7CheckRateLimitReply = mrpc.MM7CheckRateLimitReply

func (c *Client) CheckRateLimit(p MM7CheckRateLimitParams) (r *MM7CheckRateLimitReply, err error) {
	r = &MM7CheckRateLimitReply{}
	err = c.Call("CheckRateLimit", MM7CheckRateLimitParams{
		ProviderKey: p.ProviderKey,
	}, r)
	return
}
