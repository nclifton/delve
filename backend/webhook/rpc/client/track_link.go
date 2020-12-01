package client

import (
	tlrpc "github.com/burstsms/mtmo-tp/backend/webhook/rpc/types"
)

type PublishLinkHitParams = tlrpc.PublishLinkHitParams

func (c *Client) PublishLinkHit(p PublishLinkHitParams) error {
	err := c.Call("PublishLinkHit", p, &NoReply{})
	return err
}
