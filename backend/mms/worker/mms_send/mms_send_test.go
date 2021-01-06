package mms_send

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/burstsms/mtmo-tp/backend/mms/worker"
)

func TestHandle(t *testing.T) {
	t.Run("send mms to mm7 service successfully", func(t *testing.T) {
		msg := worker.Job{}
		b, _ := json.Marshal(msg)

		s := NewHandler(
			mockMM7RPCClient{},
			mockMMSRPCClient{},
		)

		err := s.Handle(context.Background(), b, map[string]interface{}{})
		if err != nil && err.Error() != "" {
			t.Error("unexpected error:", err)
		}
	})
}
