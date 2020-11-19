package client

import (
	mrpc "github.com/burstsms/mtmo-tp/backend/mm7/rpc"
)

type MM7MediaStoreParams = mrpc.MM7MediaStoreParams
type MM7MediaStoreReply = mrpc.MM7MediaStoreReply

func (c *Client) Store(p MM7MediaStoreParams) (r *MM7MediaStoreReply, err error) {
	r = &MM7MediaStoreReply{}
	err = c.Call("Store", MM7MediaStoreParams{
		FileName:    p.FileName,
		ProviderKey: p.ProviderKey,
		Extension:   p.Extension,
		Data:        p.Data,
	}, r)
	return
}
