package client

import (
	wrpc "github.com/burstsms/mtmo-tp/backend/webhook/rpc"
)

type PublishSMSStatusUpdateParams = wrpc.PublishSMSStatusUpdateParams

func (c *Client) PublishSMSStatusUpdate(params PublishSMSStatusUpdateParams) error {
	err := c.Call("PublishSMSStatusUpdate", params, &NoReply{})
	return err
}
