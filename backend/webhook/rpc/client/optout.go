package client

import (
	wrpc "github.com/burstsms/mtmo-tp/backend/webhook/rpc/types"
)

type PublishOptOutParams = wrpc.PublishOptOutParams

func (c *Client) PublishOptOut(params PublishOptOutParams) error {
	err := c.Call("PublishOptOut", params, &NoReply{})
	return err
}
