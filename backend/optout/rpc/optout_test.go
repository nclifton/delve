package rpc

import (
	"fmt"
	"testing"
)

func TestGenerateOptoutLink(t *testing.T) {

	testErr := fmt.Errorf("testerror")

	tests := []struct {
		name            string
		params          GenerateOptoutLinkParams
		db              mockDB
		expectedMessage string
		expectedErr     error
	}{
		{
			name: "test happy path with tag",
			params: GenerateOptoutLinkParams{
				AccountID:   "123",
				MessageID:   "msg_123",
				MessageType: "SMS",
				Message:     "Test message [opt-out-link]!",
			},
			db: mockDB{
				optOut: OptOut{
					LinkID: "link1",
				},
			},
			expectedMessage: "Test message http://host/link1!",
			expectedErr:     nil,
		},
		{
			name: "test happy path without OptOut tag",
			params: GenerateOptoutLinkParams{
				AccountID:   "123",
				MessageID:   "msg_123",
				MessageType: "SMS",
				Message:     "Test message!",
			},
			expectedMessage: "Test message!",
			expectedErr:     nil,
		},
		{
			name: "test with db error",
			params: GenerateOptoutLinkParams{
				AccountID:   "123",
				MessageID:   "msg_123",
				MessageType: "SMS",
				Message:     "Test message [opt-out-link]!",
			},
			db: mockDB{
				err: testErr,
			},
			expectedErr: testErr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			optOut := OptOutService{
				trackHost: "host",
				db:        test.db,
			}

			r := &GenerateOptoutLinkReply{}
			err := optOut.GenerateOptoutLink(test.params, r)
			if err != test.expectedErr {
				t.Errorf("unexpected error %+v", err)
			}

			if r.Message != test.expectedMessage {
				t.Errorf("expected Message %s, \nbut got %s", test.expectedMessage, r.Message)
			}
		})

	}
}
