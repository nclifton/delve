package fakemm7submitworker

import (
	"encoding/json"
	"errors"
	"testing"

	tcl "github.com/burstsms/mtmo-tp/backend/lib/tecloo/client"
	mm7RPC "github.com/burstsms/mtmo-tp/backend/mm7/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/mm7/worker"
)

func TestHandle(t *testing.T) {
	msg := worker.SubmitMessage{}
	b, _ := json.Marshal(msg)

	tests := map[string]struct {
		body        []byte
		rpcClient   mockRPCClient
		tecloo      mockTecloo
		expectedErr error
	}{
		"success": {
			body: b,
			rpcClient: mockRPCClient{
				rateLimitReply: mm7RPC.MM7CheckRateLimitReply{
					Allow: true,
				},
			},
			tecloo: mockTecloo{
				postMM7Response: tcl.PostMM7Response{
					Body: tcl.MM7Body{
						SubmitRsp: tcl.SubmitRsp{
							Status: tcl.Status{
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
			tecloo: mockTecloo{
				postMM7Response: tcl.PostMM7Response{
					Body: tcl.MM7Body{
						SubmitRsp: tcl.SubmitRsp{
							Status: tcl.Status{
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
		s := NewHandler(currentTest.rpcClient, currentTest.tecloo, nil)

		err := s.Handle(currentTest.body, map[string]interface{}{})
		if err != currentTest.expectedErr {
			t.Error("unexpected error:", err)
		}
	})

	t.Run("send mms to mm7 with rate limit", func(t *testing.T) {
		currentTest := tests["rate limit"]
		s := NewHandler(currentTest.rpcClient, currentTest.tecloo, nil)

		err := s.Handle(currentTest.body, map[string]interface{}{})
		if err != nil && err.Error() != currentTest.expectedErr.Error() {
			t.Error("unexpected error:", err)
		}
	})
}
