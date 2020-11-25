package client

import (
	"github.com/burstsms/mtmo-tp/backend/account/rpc"
)

type Account = rpc.Account

type FindByAPIKeyReply = rpc.FindByAPIKeyReply

func (c *Client) FindByAPIKey(key string) (r *FindByAPIKeyReply, err error) {
	r = &FindByAPIKeyReply{}
	err = c.Call("FindByAPIKey", rpc.FindByAPIKeyParams{
		Key: key,
	}, r)
	return r, err
}

type FindBySenderReply = rpc.FindBySenderReply

func (c *Client) FindBySender(sender string) (r *FindBySenderReply, err error) {
	r = &FindBySenderReply{}
	err = c.Call("FindBySender", rpc.FindBySenderParams{
		Sender: sender,
	}, r)
	return r, err
}
