package client

import (
	"encoding/gob"
	"strconv"

	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	mmsrpc "github.com/burstsms/mtmo-tp/backend/mms/rpc"
)

type NoReply = mmsrpc.NoReply

type Client struct {
	rpc.Client
}

func New(host string, port int) *Client {
	gob.Register(map[string]interface{}{})
	return &Client{
		Client: rpc.Client{
			ServiceAddress: host + ":" + strconv.Itoa(port),
			ServiceName:    mmsrpc.Name,
		},
	}
}
