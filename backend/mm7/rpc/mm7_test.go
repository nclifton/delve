package rpc

import (
	"fmt"
	"testing"
)

func TestUpdateStatus(t *testing.T) {

	testErr := fmt.Errorf("testerror")

	tests := []struct {
		name        string
		params      MM7UpdateStatusParams
		mms         mockMMS
		expectedErr error
	}{
		{
			name: "test happy path",
			params: MM7UpdateStatusParams{
				ID:          "123",
				MessageID:   "msg_123",
				Status:      "sent",
				Description: "Message sent to fake provider!",
			},
			mms:         mockMMS{},
			expectedErr: nil,
		},
		{
			name: "test with mms error",
			params: MM7UpdateStatusParams{
				ID:          "123",
				MessageID:   "msg_123",
				Status:      "sent",
				Description: "Message sent to fake provider!",
			},
			mms: mockMMS{
				err: testErr,
			},
			expectedErr: testErr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mm7 := MM7{
				svc: ConfigSvc{
					MMS: test.mms,
				},
			}

			r := &NoReply{}
			err := mm7.UpdateStatus(test.params, r)
			if err != test.expectedErr {
				t.Errorf("unexpected error %+v", err)
			}
		})

	}
}
