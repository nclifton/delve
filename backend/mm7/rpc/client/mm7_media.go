package client

import (
	"github.com/burstsms/mtmo-tp/backend/mm7/rpc/types"
)

type MediaStoreParams = types.MediaStoreParams
type MediaStoreReply = types.MediaStoreReply

func (c *Client) Store(p MediaStoreParams) (r *MediaStoreReply, err error) {
	r = &MediaStoreReply{}
	err = c.Call("Store", MediaStoreParams{
		FileName:    p.FileName,
		ProviderKey: p.ProviderKey,
		Extension:   p.Extension,
		Data:        p.Data,
	}, r)
	return
}
