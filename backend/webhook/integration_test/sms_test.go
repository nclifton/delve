// +build integration

package test

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burstsms/mtmo-tp/backend/lib/assertdb"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/service"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
)

type ExpectedSMSStatusData = service.PublishStatusData

type ExpectedSMSStatusRequestBody struct {
	Event string                `json:"event"`
	Data  ExpectedSMSStatusData `json:"data"`
}

func Test_PublishSMSStatusUpdate(t *testing.T) {
	log.Println("test PublishSMSStatusUpdate")

	i := setupForPublishSMSStatusUpdate(t)
	defer i.teardown(t)
	client := i.getClient(t)
	timestampNow := timestamppb.Now()

	type wantErr struct {
		status *status.Status
		ok     bool
	}

	tests := []struct {
		name    string
		params  *webhookpb.PublishSMSStatusUpdateParams
		want    ExpectedRequests
		wantErr wantErr
	}{
		{
			name: "happy",
			params: &webhookpb.PublishSMSStatusUpdateParams{
				AccountId:       "42",
				SMSId:           "xxy",
				MessageRef:      "123",
				Recipient:       "35426378914",
				Sender:          "46354078643",
				Status:          "done",
				StatusUpdatedAt: timestampNow,
			},
			want: ExpectedRequests{
				NumberOfRequests: 1,
				WaitMilliseconds: 500,
				Methods:          []string{"POST"},
				ContentTypes:     []string{"application/json"},
				Bodies: []string{
					jsonString(t,
						ExpectedSMSStatusRequestBody{
							service.EventSMSStatus,
							ExpectedSMSStatusData{
								SMS_id:            "xxy",
								Message_ref:       "123",
								Recipient:         "35426378914",
								Sender:            "46354078643",
								Status:            "done",
								Status_updated_at: timestampNow.AsTime().Format(time.RFC3339),
							}})}},
			wantErr: wantErr{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("test: %s", tt.name)
			i.ResetHttpRequests()
			got, err := client.PublishSMSStatusUpdate(i.ctx, tt.params)
			if tt.wantErr.status != nil && err != nil {
				errStatus, ok := status.FromError(err)
				assert.Equal(t, ok, tt.wantErr.ok, "grpc ok")
				assert.EqualValues(t, tt.wantErr.status, errStatus, "grpc status")
			} else if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			i.WaitForRequests(tt.want)
			assert.ObjectsAreEqual(webhookpb.NoReply{}, got)
			i.AssertRequests(tt.want)
		})
	}

}

func setupForPublishSMSStatusUpdate(t *testing.T) *testDeps {
	i := newSetup(t, tfx)
	i.HaveInDatabase("webhook", assertdb.Row{
		"id":         32767,
		"account_id": "42",
		"event":      service.EventSMSStatus,
		"name":       "name1",
		"url":        i.webhookURL,
		"rate_limit": 2,
		"created_at": "2020-01-12 22:41:42",
		"updated_at": "2020-01-12 22:41:42"})
	return i
}

type ExpectedLastMessage = service.PublishMessageData
type ExpectedMOStatusData = service.PublishMOData
type ExpectedMOStatusRequestBody struct {
	Event string               `json:"event"`
	Data  ExpectedMOStatusData `json:"data"`
}

func Test_PublishMO(t *testing.T) {
	log.Println("test PublishOptOut")

	i := setupForPublishMO(t)
	defer i.teardown(t)
	client := i.getClient(t)
	timestampNow := timestamppb.Now()

	type wantErr struct {
		status *status.Status
		ok     bool
	}
	type want struct {
		reply    *webhookpb.NoReply
		requests ExpectedRequests
	}

	tests := []struct {
		name    string
		params  *webhookpb.PublishMOParams
		want    want
		wantErr wantErr
	}{
		{
			name: "happy",
			params: &webhookpb.PublishMOParams{
				AccountId:  "42",
				SMSId:      "xxy",
				Message:    "General Kenobi",
				Recipient:  "35426378914",
				Sender:     "46354078643",
				ReceivedAt: timestampNow,
				LastMessage: &webhookpb.Message{
					Type:        "sms",
					Id:          "21",
					Recipient:   "46354078643",
					Sender:      "35426378914",
					Message:     "Hello there",
					MessageRef:  "abc",
					Subject:     "Greetings",
					ContentURLs: []string{"http://example.com/dickpic.png"},
				},
			},
			want: want{
				reply: &webhookpb.NoReply{},
				requests: ExpectedRequests{
					NumberOfRequests: 1,
					WaitMilliseconds: 500,
					Methods:          []string{"POST"},
					ContentTypes:     []string{"application/json"},
					Bodies: []string{
						jsonString(t,
							ExpectedMOStatusRequestBody{
								service.EventMOStatus,
								ExpectedMOStatusData{
									SMS_id:    "xxy",
									Recipient: "35426378914",
									Sender:    "46354078643",
									Message:   "General Kenobi",
									Timestamp: timestampNow.AsTime().Format(time.RFC3339),
									Last_message: ExpectedLastMessage{
										Type:         "sms",
										Id:           "21",
										Recipient:    "46354078643",
										Sender:       "35426378914",
										Subject:      "Greetings",
										Message:      "Hello there",
										Content_urls: []string{"http://example.com/dickpic.png"},
										Message_ref:  "abc",
									}}})}}},
			wantErr: wantErr{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("test: %s", tt.name)
			i.ResetHttpRequests()
			got, err := client.PublishMO(i.ctx, tt.params)
			if tt.wantErr.status != nil && err != nil {
				errStatus, ok := status.FromError(err)
				assert.Equal(t, ok, tt.wantErr.ok, "grpc ok")
				assert.EqualValues(t, tt.wantErr.status, errStatus, "grpc status")
			} else if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			assert.ObjectsAreEqual(tt.want.reply, got)
			i.WaitForRequests(tt.want.requests)
			i.AssertRequests(tt.want.requests)
		})
	}

}

func setupForPublishMO(t *testing.T) *testDeps {
	i := newSetup(t, tfx)
	i.HaveInDatabase("webhook", assertdb.Row{
		"id":         32767,
		"account_id": "42",
		"event":      service.EventMOStatus,
		"name":       "name1",
		"url":        i.webhookURL,
		"rate_limit": 2,
		"created_at": "2020-01-12 22:41:42",
		"updated_at": "2020-01-12 22:41:42"})
	return i
}
