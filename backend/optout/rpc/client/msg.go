package client

import (
	"github.com/burstsms/mtmo-tp/backend/optout/rpc/types"
)

type OptOutViaMsgParams = types.OptOutViaMsgParams

func (c *Client) OptOutViaMsg(params OptOutViaMsgParams) error {
	return c.Call("OptOutViaMsg", params, &types.NoReply{})
}
