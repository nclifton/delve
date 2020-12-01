//+build integration

package client_test

import (
	"log"
	"testing"

	mm7RPC "github.com/burstsms/mtmo-tp/backend/mm7/rpc"
	"github.com/burstsms/mtmo-tp/backend/mm7/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/mm7/rpc/types"
	"github.com/kelseyhightower/envconfig"
)

func initClient() *client.Client {
	var env mm7RPC.Env
	err := envconfig.Process("mm7", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	return client.NewClient(env.RPCHost, env.RPCPort)
}
func TestPing(t *testing.T) {
	cli := initClient()

	r, err := cli.Ping()
	if err != nil {
		log.Fatal("failed to call ping:", err)
	}

	if r.Res != "PONG" {
		log.Fatal("failed wrong response:", r)
	}
}

func TestSend(t *testing.T) {
	cli := initClient()

	tests := []struct {
		name        string
		sendParams  types.SendParams
		expectedErr error
	}{
		{
			name: "test fakesubmit happy path",
			sendParams: types.SendParams{
				ID:          "msg id",
				Subject:     "subject",
				Message:     "message",
				Sender:      "sender",
				Recipient:   "recipient",
				ContentURLs: []string{"https://res.cloudinary.com/burstsms/image/upload/v1550725703/bqe9vphnyu0scdblpf08.jpg"},
				ProviderKey: "fake",
			},
			expectedErr: nil,
		},
		{
			name: "test fakesubmit happy path with client error",
			sendParams: types.SendParams{
				ID:          "msg id",
				Subject:     "subject",
				Message:     "message",
				Sender:      "sender",
				Recipient:   "61422262000",
				ContentURLs: []string{"https://res.cloudinary.com/burstsms/image/upload/v1550725703/bqe9vphnyu0scdblpf08.jpg"},
				ProviderKey: "fake",
			},
			expectedErr: nil,
		},
		{
			name: "test otpussubmit happy path",
			sendParams: types.SendParams{
				ID:          "msg id",
				Subject:     "subject",
				Message:     "message",
				Sender:      "sender",
				Recipient:   "recipient",
				ContentURLs: []string{"https://res.cloudinary.com/burstsms/image/upload/v1550725703/bqe9vphnyu0scdblpf08.jpg"},
				ProviderKey: "optus",
			},
			expectedErr: nil,
		},
		{
			name: "test mgagesubmit happy path",
			sendParams: types.SendParams{
				ID:          "msg id",
				Subject:     "subject",
				Message:     "message",
				Sender:      "sender",
				Recipient:   "recipient",
				ContentURLs: []string{"https://res.cloudinary.com/burstsms/image/upload/v1550725703/bqe9vphnyu0scdblpf08.jpg"},
				ProviderKey: "mgage",
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		_, err := cli.Send(test.sendParams)

		if err != test.expectedErr {
			t.Errorf("for %s, \nexpected error %v, \n but got error %v", test.name, test.expectedErr, err)
		}
	}
}
