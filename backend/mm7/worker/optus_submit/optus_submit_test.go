package optussubmitworker

import (
	"encoding/json"
	"errors"
	"testing"

	optcl "github.com/burstsms/mtmo-tp/backend/lib/optus/client"
	mm7RPC "github.com/burstsms/mtmo-tp/backend/mm7/rpc"
	"github.com/burstsms/mtmo-tp/backend/mm7/worker"
)

func TestHandle(t *testing.T) {
	msg := worker.SubmitMessage{}
	b, _ := json.Marshal(msg)

	tests := map[string]struct {
		body        []byte
		rpcClient   mockRPCClient
		optus       mockOptus
		expectedErr error
	}{
		"success": {
			body: b,
			rpcClient: mockRPCClient{
				rateLimitReply: mm7RPC.MM7CheckRateLimitReply{
					Allow: true,
				},
			},
			optus: mockOptus{
				postMM7Response: optcl.PostMM7Response{
					Body: optcl.MM7Body{
						SubmitRsp: optcl.SubmitRsp{
							Status: optcl.Status{
								StatusCode: "1000",
								StatusText: "success",
							},
						},
					},
				},
				statusCode: 200,
			}, expectedErr: nil,
		},
		"rate limit": {
			body: b,
			rpcClient: mockRPCClient{
				rateLimitReply: mm7RPC.MM7CheckRateLimitReply{
					Allow: false,
				},
			},
			optus: mockOptus{
				postMM7Response: optcl.PostMM7Response{
					Body: optcl.MM7Body{
						SubmitRsp: optcl.SubmitRsp{
							Status: optcl.Status{
								StatusCode: "1000",
								StatusText: "success",
							},
						},
					},
				},
				statusCode: 200,
			}, expectedErr: errors.New("Failed sending message id:  Error: rate limit reached"),
		},
	}

	t.Run("send mms to mm7 successfully", func(t *testing.T) {
		currentTest := tests["success"]
		s := NewHandler(currentTest.rpcClient, currentTest.optus, nil)

		err := s.Handle(currentTest.body, map[string]interface{}{})
		if err != currentTest.expectedErr {
			t.Error("unexpected error:", err)
		}
	})

	t.Run("send mms to mm7 with rate limit", func(t *testing.T) {
		currentTest := tests["rate limit"]
		s := NewHandler(currentTest.rpcClient, currentTest.optus, nil)

		err := s.Handle(currentTest.body, map[string]interface{}{})
		if err != nil && err.Error() != currentTest.expectedErr.Error() {
			t.Error("unexpected error:", err)
		}
	})
}
