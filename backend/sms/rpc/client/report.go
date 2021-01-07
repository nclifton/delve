package client

import (
	rpc "github.com/burstsms/mtmo-tp/backend/sms/rpc/types"
)

type GenerateAccountSMSUsageParams = rpc.GenerateAccountSMSUsageParams
type GenerateAccountSMSUsageResponse = rpc.GenerateAccountSMSUsageResponse

func (c *Client) GenerateAccountSMSUsage(p GenerateAccountSMSUsageParams) (r *GenerateAccountSMSUsageResponse, err error) {
	r = &GenerateAccountSMSUsageResponse{}
	err = c.Call("GenerateAccountSMSUsage", p, r)
	return r, err
}
