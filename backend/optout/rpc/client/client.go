package client

import (
	"encoding/gob"
	"strconv"

	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	"github.com/burstsms/mtmo-tp/backend/optout/rpc/types"
)

type NoReply = types.NoReply

type Client struct {
	rpc.Client
}

func NewClient(host string, port int) *Client {
	gob.Register(map[string]interface{}{})
	return &Client{
		Client: rpc.Client{
			ServiceAddress: host + ":" + strconv.Itoa(port),
			ServiceName:    types.Name,
		},
	}
}
