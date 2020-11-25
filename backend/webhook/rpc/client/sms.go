package client

import (
	wrpc "github.com/burstsms/mtmo-tp/backend/webhook/rpc"
)

type PublishSMSStatusUpdateParams = wrpc.PublishSMSStatusUpdateParams

func (c *Client) PublishSMSStatusUpdate(params PublishSMSStatusUpdateParams) error {
	return c.Call("PublishSMSStatusUpdate", params, &NoReply{})
}

type PublishMOParams = wrpc.PublishMOParams
type LastMessage = wrpc.LastMessage

func (c *Client) PublishMO(params PublishMOParams) error {
	err := c.Call("PublishMO", params, &NoReply{})
	return err
}
