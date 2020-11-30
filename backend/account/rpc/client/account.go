package client

import (
	"github.com/burstsms/mtmo-tp/backend/account/rpc/types"
)

type Account = types.Account

type FindByAPIKeyReply = types.FindByAPIKeyReply

func (c *Client) FindByAPIKey(key string) (r *FindByAPIKeyReply, err error) {
	r = &FindByAPIKeyReply{}
	err = c.Call("FindByAPIKey", types.FindByAPIKeyParams{
		Key: key,
	}, r)
	return r, err
}

type FindBySenderReply = types.FindBySenderReply

func (c *Client) FindBySender(sender string) (r *FindBySenderReply, err error) {
	r = &FindBySenderReply{}
	err = c.Call("FindBySender", types.FindBySenderParams{
		Sender: sender,
	}, r)
	return r, err
}
