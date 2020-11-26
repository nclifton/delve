package client

import (
	"encoding/gob"
	"strconv"

	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	orpc "github.com/burstsms/mtmo-tp/backend/optout/rpc"
)

type NoReply = orpc.NoReply

type Client struct {
	rpc.Client
}

func New(host string, port int) *Client {
	gob.Register(map[string]interface{}{})
	return &Client{
		Client: rpc.Client{
			ServiceAddress: host + ":" + strconv.Itoa(port),
			ServiceName:    orpc.Name,
		},
	}
}
