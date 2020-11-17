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
