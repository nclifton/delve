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
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/app/service"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
)

type ExpectedLinkHitMessage = service.PublishMessageData
type ExpectedLinkHitData = service.PublishLinkHitData
type ExpectedLinkHitRequestBody struct {
	Event string              `json:"event"`
	Data  ExpectedLinkHitData `json:"data"`
}

func Test_PublishLinkHit(t *testing.T) {
	log.Println("test PublishOptOut")

	i := setupForPublishLinkHit(t)
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
		params  *webhookpb.PublishLinkHitParams
		want    want
		wantErr wantErr
	}{
		{
			name: "happy",
			params: &webhookpb.PublishLinkHitParams{
				URL:       "http://example.com/hitme",
				Hits:      1,
				Timestamp: timestampNow,
				SourceMessage: &webhookpb.Message{
					Type:       "sms",
					Id:         "21",
					Recipient:  "46354078643",
					Sender:     "35426378914",
					Message:    "General Kenobi",
					MessageRef: "abc",
				},
				AccountId: "42",
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
							ExpectedLinkHitRequestBody{
								service.EventLinkHitStatus,
								ExpectedLinkHitData{
									URL:       "http://example.com/hitme",
									Hits:      1,
									Timestamp: timestampNow.AsTime().Format(time.RFC3339),
									Source_message: service.PublishMessageData{
										Type:        "sms",
										Id:          "21",
										Recipient:   "46354078643",
										Sender:      "35426378914",
										Message:     "General Kenobi",
										Message_ref: "abc",
									}}})}}},
			wantErr: wantErr{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("test: %s", tt.name)
			i.ResetHttpRequests()
			got, err := client.PublishLinkHit(i.ctx, tt.params)
			if tt.wantErr.status != nil && err != nil {
				errStatus, ok := status.FromError(err)
				assert.Equal(t, ok, tt.wantErr.ok, "grpc ok")
				assert.EqualValues(t, tt.wantErr.status, errStatus, "grpc status")
			} else if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			assert.ObjectsAreEqual(webhookpb.NoReply{}, got)
			i.WaitForRequests(tt.want.requests)
			i.AssertRequests(tt.want.requests)
		})
	}

}

func setupForPublishLinkHit(t *testing.T) *testDeps {
	s := newSetup(t, tfx)
	s.HaveInDatabase("webhook", assertdb.Row{
		"id":         32767,
		"account_id": "42",
		"event":      service.EventLinkHitStatus,
		"name":       "name1",
		"url":        s.webhookURL,
		"rate_limit": 2,
		"created_at": "2020-01-12 22:41:42",
		"updated_at": "2020-01-12 22:41:42"})
	return s
}
