package client

import (
	"encoding/gob"
	"strconv"

	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	tlrpc "github.com/burstsms/mtmo-tp/backend/track_link/rpc"
)

type NoReply = tlrpc.NoReply

type Client struct {
	rpc.Client
}

func NewClient(host string, port int) *Client {
	gob.Register(map[string]interface{}{})
	return &Client{
		Client: rpc.Client{
			ServiceAddress: host + ":" + strconv.Itoa(port),
			ServiceName:    tlrpc.Name,
		},
	}
}
