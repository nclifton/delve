// +build integration

package test

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/app/service"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
)

type ExpectedSMSStatusData = service.PublishStatusData

type ExpectedSMSStatusRequestBody struct {
	Event string                `json:"event"`
	Data  ExpectedSMSStatusData `json:"data"`
}

func Test_PublishSMSStatusUpdate(t *testing.T) {
	log.Println("test PublishSMSStatusUpdate")

	setup := setupForPublishSMSStatusUpdate(t)
	defer setup.teardown(t)
	client := setup.getClient(t)
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
				numberOfRequests: 1,
				waitMilliseconds: 500,
				methods:          []string{"POST"},
				contentTypes:     []string{"application/json"},
				bodies: []string{
					setup.marshalJson(t,
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
			setup.resetHttpRequests(t)
			got, err := client.PublishSMSStatusUpdate(setup.ctx, tt.params)
			if tt.wantErr.status != nil && err != nil {
				errStatus, ok := status.FromError(err)
				assert.Equal(t, ok, tt.wantErr.ok, "grpc ok")
				assert.EqualValues(t, tt.wantErr.status, errStatus, "grpc status")
			} else if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			setup.waitForRequests(t, tt.want)
			assert.ObjectsAreEqual(webhookpb.NoReply{}, got)
			setup.assertRequests(t, tt.want)
		})
	}

}

func setupForPublishSMSStatusUpdate(t *testing.T) *testDeps {
	setup := newSetup(t, tfx, listener)
	setup.startHttpServer(t)
	setup.startWorker(t)
	setup.adb.HaveInDatabase("webhook",
		"id, account_id, event, name, url, rate_limit, created_at, updated_at",
		[]interface{}{32767, "42", service.EventSMSStatus, "name1", setup.httpServer.URL, 2, "2020-01-12 22:41:42", "2020-01-12 22:41:42"})
	return setup
}

type ExpectedLastMessage = service.PublishMessageData
type ExpectedMOStatusData = service.PublishMOData
type ExpectedMOStatusRequestBody struct {
	Event string               `json:"event"`
	Data  ExpectedMOStatusData `json:"data"`
}

func Test_PublishMO(t *testing.T) {
	log.Println("test PublishOptOut")

	setup := setupForPublishMO(t)
	defer setup.teardown(t)
	client := setup.getClient(t)
	timestampNow := timestamppb.Now()

	type wantErr struct {
		status *status.Status
		ok     bool
	}

	tests := []struct {
		name    string
		params  *webhookpb.PublishMOParams
		want    ExpectedRequests
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
			want: ExpectedRequests{
				numberOfRequests: 1,
				waitMilliseconds: 500,
				methods:          []string{"POST"},
				contentTypes:     []string{"application/json"},
				bodies: []string{
					setup.marshalJson(t,
						ExpectedMOStatusRequestBody{
							service.EventMOStatus,
							ExpectedMOStatusData{
								SMS_id:       "xxy",
								Recipient:    "35426378914",
								Sender:       "46354078643",
								Message:      "General Kenobi",
								Timestamp:    timestampNow.AsTime().Format(time.RFC3339),
								Last_message: ExpectedLastMessage{
									Type:         "sms",
									Id:           "21",
									Recipient:    "46354078643",
									Sender:       "35426378914",
									Subject:      "Greetings",
									Message:      "Hello there",
									Content_urls: []string{"http://example.com/dickpic.png"},
									Message_ref:  "abc",
								},
							}})}},
			wantErr: wantErr{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("test: %s", tt.name)
			setup.resetHttpRequests(t)
			got, err := client.PublishMO(setup.ctx, tt.params)
			if tt.wantErr.status != nil && err != nil {
				errStatus, ok := status.FromError(err)
				assert.Equal(t, ok, tt.wantErr.ok, "grpc ok")
				assert.EqualValues(t, tt.wantErr.status, errStatus, "grpc status")
			} else if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			setup.waitForRequests(t, tt.want)
			assert.ObjectsAreEqual(webhookpb.NoReply{}, got)
			setup.assertRequests(t, tt.want)
		})
	}

}

func setupForPublishMO(t *testing.T) *testDeps {
	setup := newSetup(t, tfx, listener)
	setup.startHttpServer(t)
	setup.startWorker(t)
	setup.adb.HaveInDatabase("webhook",
		"id, account_id, event, name, url, rate_limit, created_at, updated_at",
		[]interface{}{32767, "42", service.EventMOStatus, "name1", setup.httpServer.URL, 2, "2020-01-12 22:41:42", "2020-01-12 22:41:42"})
	return setup
}
