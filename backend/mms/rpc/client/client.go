package client

import (
	"encoding/gob"

	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	"github.com/burstsms/mtmo-tp/backend/mms/rpc/types"
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
