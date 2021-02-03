package rpc

import (
	"fmt"
	"testing"

	mms "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
	types "github.com/burstsms/mtmo-tp/backend/optout/rpc/types"
	sms "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
)

func TestFindByLinkID(t *testing.T) {

	testErr := fmt.Errorf("testerror")

	tests := []struct {
		name          string
		params        types.FindByLinkIDParams
		db            mockDB
		expectedReply types.FindByLinkIDReply
		expectedErr   error
	}{
		{
			name: "test happy path",
			params: types.FindByLinkIDParams{
				LinkID: "dVYHEhq6",
			},
			db: mockDB{
				optOut: types.OptOut{
					ID:          "11111111-1111-1111-1111-111111111111",
					AccountID:   "11111111-1111-1111-1111-111111111112",
					MessageID:   "11111111-1111-1111-1111-111111111113",
					MessageType: "SMS",
					LinkID:      "dVYHEhq6",
				},
			},
			expectedReply: types.FindByLinkIDReply{
				OptOut: &types.OptOut{
					ID:          "11111111-1111-1111-1111-111111111111",
					AccountID:   "11111111-1111-1111-1111-111111111112",
					MessageID:   "11111111-1111-1111-1111-111111111113",
					MessageType: "SMS",
					LinkID:      "dVYHEhq6",
				},
			},
			expectedErr: nil,
		},
		{
			name: "test with db error",
			params: types.FindByLinkIDParams{
				LinkID: "dVYHEhq6",
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
				db: test.db,
			}

			r := &types.FindByLinkIDReply{}
			err := optOut.FindByLinkID(test.params, r)
			if err != test.expectedErr {
				t.Errorf("unexpected error %+v", err)
			}

			if err == nil && !types.OptOutEqual(*test.expectedReply.OptOut, *r.OptOut) {
				t.Errorf("unexpected result %+v, \nbut got %+v", test.expectedReply.OptOut, r.OptOut)
			}

		})
	}
}

func TestOptOutViaLink(t *testing.T) {

	testErr := fmt.Errorf("testerror")

	tests := []struct {
		name          string
		params        types.OptOutViaLinkParams
		db            mockDB
		webhookRPC    mockWebhookRPC
		smsRPC        mockSMSRPC
		mmsRPC        mockMMSRPC
		expectedReply types.OptOutViaLinkReply
		expectedErr   error
	}{
		{
			name: "test happy path sms with tag",
			params: types.OptOutViaLinkParams{
				LinkID: "dVYHEhq6",
			},
			db: mockDB{
				optOut: types.OptOut{
					ID:          "11111111-1111-1111-1111-111111111111",
					AccountID:   "11111111-1111-1111-1111-111111111112",
					MessageID:   "11111111-1111-1111-1111-111111111113",
					MessageType: "sms",
					LinkID:      "dVYHEhq6",
				},
			},
			smsRPC: mockSMSRPC{
				findByIDReply: sms.FindByIDReply{
					SMS: &sms.SMS{},
				},
			},
			expectedReply: types.OptOutViaLinkReply{
				OptOut: &types.OptOut{
					ID:          "11111111-1111-1111-1111-111111111111",
					AccountID:   "11111111-1111-1111-1111-111111111112",
					MessageID:   "11111111-1111-1111-1111-111111111113",
					MessageType: "sms",
					LinkID:      "dVYHEhq6",
				},
			},
			expectedErr: nil,
		},
		{
			name: "test happy path mms",
			params: types.OptOutViaLinkParams{
				LinkID: "dVYHEhq6",
			},
			db: mockDB{
				optOut: types.OptOut{
					ID:          "11111111-1111-1111-1111-111111111111",
					AccountID:   "11111111-1111-1111-1111-111111111112",
					MessageID:   "11111111-1111-1111-1111-111111111113",
					MessageType: "mms",
					LinkID:      "dVYHEhq6",
				},
			},
			mmsRPC: mockMMSRPC{
				findByIDReply: mms.FindByIDReply{
					MMS: &mms.MMS{},
				},
			},
			expectedReply: types.OptOutViaLinkReply{
				OptOut: &types.OptOut{
					ID:          "11111111-1111-1111-1111-111111111111",
					AccountID:   "11111111-1111-1111-1111-111111111112",
					MessageID:   "11111111-1111-1111-1111-111111111113",
					MessageType: "mms",
					LinkID:      "dVYHEhq6",
				},
			},
			expectedErr: nil,
		},
		{
			name: "test with db error",
			params: types.OptOutViaLinkParams{
				LinkID: "dVYHEhq6",
			},
			db: mockDB{
				err: testErr,
			},
			expectedErr: testErr,
		},
		{
			name: "test with webhook rpc error",
			params: types.OptOutViaLinkParams{
				LinkID: "dVYHEhq6",
			},
			db: mockDB{
				optOut: types.OptOut{
					ID:          "11111111-1111-1111-1111-111111111111",
					AccountID:   "11111111-1111-1111-1111-111111111112",
					MessageID:   "11111111-1111-1111-1111-111111111113",
					MessageType: "sms",
					LinkID:      "dVYHEhq6",
				},
			},
			smsRPC: mockSMSRPC{
				findByIDReply: sms.FindByIDReply{
					SMS: &sms.SMS{},
				},
			},
			webhookRPC: mockWebhookRPC{
				err: testErr,
			},
			expectedErr: testErr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			optOut := OptOutService{
				db:         test.db,
				webhookRPC: test.webhookRPC,
				smsRPC:     test.smsRPC,
				mmsRPC:     test.mmsRPC,
			}

			r := &types.OptOutViaLinkReply{}
			err := optOut.OptOutViaLink(test.params, r)
			if err != test.expectedErr {
				t.Errorf("unexpected error %+v", err)
			}

			if err == nil && !types.OptOutEqual(*test.expectedReply.OptOut, *r.OptOut) {
				t.Errorf("unexpected result %+v, \nbut got %+v", test.expectedReply.OptOut, r.OptOut)
			}

		})
	}
}

func TestGenerateOptOutLink(t *testing.T) {

	testErr := fmt.Errorf("testerror")

	tests := []struct {
		name            string
		params          types.GenerateOptOutLinkParams
		db              mockDB
		expectedMessage string
		expectedErr     error
	}{
		{
			name: "test happy path with tag",
			params: types.GenerateOptOutLinkParams{
				AccountID:   "123",
				MessageID:   "msg_123",
				MessageType: "SMS",
				Message:     "Test message [opt-out-link]!",
				Sender:      "1701",
			},
			db: mockDB{
				optOut: types.OptOut{
					LinkID: "link1",
				},
			},
			expectedMessage: "Test message http://host/link1!",
			expectedErr:     nil,
		},
		{
			name: "test happy path without OptOut tag",
			params: types.GenerateOptOutLinkParams{
				AccountID:   "123",
				MessageID:   "msg_123",
				MessageType: "SMS",
				Message:     "Test message!",
				Sender:      "1701",
			},
			expectedMessage: "Test message!",
			expectedErr:     nil,
		},
		{
			name: "test with db error",
			params: types.GenerateOptOutLinkParams{
				AccountID:   "123",
				MessageID:   "msg_123",
				MessageType: "SMS",
				Message:     "Test message [opt-out-link]!",
				Sender:      "1701",
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
				optOutDomain: "host",
				db:           test.db,
			}

			r := &types.GenerateOptOutLinkReply{}
			err := optOut.GenerateOptOutLink(test.params, r)
			if err != test.expectedErr {
				t.Errorf("unexpected error %+v", err)
			}

			if r.Message != test.expectedMessage {
				t.Errorf("expected Message %s, \nbut got %s", test.expectedMessage, r.Message)
			}
		})
	}
}
