// +build integration

package test

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burstsms/mtmo-tp/backend/lib/asserthttp"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/app/service"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
)

type ExpectedOptOutSourceMessage = service.PublishMessageData
type ExpectedOptOutData = service.PublishOptOutData
type ExpectedOptOutRequestBody struct {
	Event string             `json:"event"`
	Data  ExpectedOptOutData `json:"data"`
}

func Test_PublishOptOut(t *testing.T) {
	log.Println("test PublishOptOut")

	i := setupForPublishOptOut(t)
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
		params  *webhookpb.PublishOptOutParams
		want    want
		wantErr wantErr
	}{
		{
			name: "happy opt out via sms",
			params: &webhookpb.PublishOptOutParams{
				Source:    "here",
				Timestamp: timestampNow,
				SourceMessage: &webhookpb.Message{
					Type:        "sms",
					Id:          "xxy",
					Recipient:   "35426378914",
					Sender:      "46354078643",
					Message:     "Hello there",
					MessageRef:  "123",
					Subject:     "Greetings",
					ContentURLs: []string{"http://example.com/image.png"},
				},
				AccountId: "42",
			},
			want: want{
				reply: &webhookpb.NoReply{},
				requests: asserthttp.ExpectedRequests{
					NumberOfRequests: 1,
					WaitMilliseconds: 500,
					Methods:          []string{"POST"},
					ContentTypes:     []string{"application/json"},
					Bodies: []string{
						jsonString(t,
							ExpectedOptOutRequestBody{
								service.EventOptOutStatus,
								ExpectedOptOutData{
									Source:    "here",
									Timestamp: timestampNow.AsTime().Format(time.RFC3339),
									Source_message: ExpectedOptOutSourceMessage{
										Type:        "sms",
										Id:          "xxy",
										Recipient:   "35426378914",
										Sender:      "46354078643",
										Message:     "Hello there",
										Message_ref: "123",
									}}})}}},
			wantErr: wantErr{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("test: %s", tt.name)
			i.ResetHttpRequests()
			got, err := client.PublishOptOut(i.ctx, tt.params)
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

func setupForPublishOptOut(t *testing.T) *testDeps {
	i := newSetup(t, tfx)
	i.HaveInDatabase("webhook",
		"id, account_id, event, name, url, rate_limit, created_at, updated_at",
		[]interface{}{32767, "42", service.EventOptOutStatus, "name1", i.webhookURL, 2, "2020-01-12 22:41:42", "2020-01-12 22:41:42"})
	return i
}
