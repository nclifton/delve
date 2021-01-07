package client

import (
	rpc "github.com/burstsms/mtmo-tp/backend/mms/rpc/types"
)

type GenerateAccountMMSUsageParams = rpc.GenerateAccountMMSUsageParams
type GenerateAccountMMSUsageResponse = rpc.GenerateAccountMMSUsageResponse

func (c *Client) GenerateAccountMMSUsage(p GenerateAccountMMSUsageParams) (r *GenerateAccountMMSUsageResponse, err error) {
	r = &GenerateAccountMMSUsageResponse{}
	err = c.Call("GenerateAccountMMSUsage", p, r)
	return r, err
}
