package client

import "github.com/burstsms/mtmo-tp/backend/webhook/rpc/types"

type PublishSMSStatusUpdateParams = types.PublishSMSStatusUpdateParams

func (c *Client) PublishSMSStatusUpdate(params PublishSMSStatusUpdateParams) error {
	return c.Call("PublishSMSStatusUpdate", params, &types.NoReply{})
}

type PublishMOParams = types.PublishMOParams
type LastMessage = types.LastMessage

func (c *Client) PublishMO(params PublishMOParams) error {
	return c.Call("PublishMO", params, &types.NoReply{})
}
