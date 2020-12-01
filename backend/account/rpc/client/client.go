package client

import (
	"encoding/gob"
	"strconv"

	"github.com/burstsms/mtmo-tp/backend/account/rpc/types"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
)

type Client struct {
	rpc.Client
}

func New(host string, port int) *Client {
	gob.Register(map[string]interface{}{})
	return &Client{
		Client: rpc.Client{
			ServiceAddress: host + ":" + strconv.Itoa(port),
			ServiceName:    types.Name,
		},
	}
}
