package client

import (
	wrpc "github.com/burstsms/mtmo-tp/backend/webhook/rpc/types"
)

type PublishMMSStatusUpdateParams = wrpc.PublishMMSStatusUpdateParams

func (c *Client) PublishMMSStatusUpdate(params PublishMMSStatusUpdateParams) error {
	return c.Call("PublishMMSStatusUpdate", params, &NoReply{})
}
