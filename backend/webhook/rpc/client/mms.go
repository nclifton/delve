package client

import (
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/types"
)

type PublishMMSStatusUpdateParams = types.PublishMMSStatusUpdateParams

func (c *Client) PublishMMSStatusUpdate(p PublishMMSStatusUpdateParams) error {
	return c.Call("PublishMMSStatusUpdate", p, &types.NoReply{})
}
