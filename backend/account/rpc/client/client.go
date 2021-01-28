package client

import (
	"encoding/gob"

	"github.com/burstsms/mtmo-tp/backend/account/rpc/types"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
)

type Client struct {
	rpc.Client
}

func New(addr string) *Client {
	gob.Register(map[string]interface{}{})
	return &Client{
		Client: rpc.Client{
			ServiceAddress: addr,
			ServiceName:    types.Name,
		},
	}
}
