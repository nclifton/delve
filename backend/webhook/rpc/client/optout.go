package client

import "github.com/burstsms/mtmo-tp/backend/webhook/rpc/types"

type PublishOptOutParams = types.PublishOptOutParams

func (c *Client) PublishOptOut(p PublishOptOutParams) error {
	return c.Call("PublishOptOut", p, &types.NoReply{})
}
