package client

import (
	"encoding/gob"
	"strconv"

	arpc "github.com/burstsms/mtmo-tp/backend/account/rpc"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
)

type NoReply = arpc.NoReply

type Client struct {
	rpc.Client
}

func New(host string, port int) *Client {
	gob.Register(map[string]interface{}{})
	return &Client{
		Client: rpc.Client{
			ServiceAddress: host + ":" + strconv.Itoa(port),
			ServiceName:    arpc.Name,
		},
	}
}
