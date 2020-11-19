package client

import (
	"encoding/gob"
	"strconv"

	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	mrpc "github.com/burstsms/mtmo-tp/backend/mm7/rpc"
)

type NoReply = mrpc.NoReply

type Client struct {
	rpc.Client
}

func NewClient(host string, port int) *Client {
	gob.Register(map[string]interface{}{})
	return &Client{
		Client: rpc.Client{
			ServiceAddress: host + ":" + strconv.Itoa(port),
			ServiceName:    mrpc.Name,
		},
	}
}
