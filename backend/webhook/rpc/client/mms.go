package client

import (
	wrpc "github.com/burstsms/mtmo-tp/backend/webhook/rpc"
)

type PublishMMSStatusUpdateParams = wrpc.PublishMMSStatusUpdateParams

func (c *Client) PublishMMSStatusUpdate(params PublishMMSStatusUpdateParams) error {
	return c.Call("PublishMMSStatusUpdate", params, &NoReply{})
}
