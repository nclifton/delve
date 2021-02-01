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

type FindByIDReply = types.FindByIDReply

func (c *Client) FindByID(id string) (r *FindByIDReply, err error) {
	r = &FindByIDReply{}
	err = c.Call("FindByID", types.FindByIDParams{
		ID: id,
	}, r)
	return r, err
}
