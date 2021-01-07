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

type ExpectedLinkHitMessage = service.PublishMessageData
type ExpectedLinkHitData = service.PublishLinkHitData
type ExpectedLinkHitRequestBody struct {
	Event string              `json:"event"`
	Data  ExpectedLinkHitData `json:"data"`
}

func Test_PublishLinkHit(t *testing.T) {
	log.Println("test PublishOptOut")

	setup := setupForPublishLinkHit(t)
	defer setup.teardown(t)
	client := setup.getClient(t)
	timestampNow := timestamppb.Now()

	type wantErr struct {
		status *status.Status
		ok     bool
	}

	tests := []struct {
		name    string
		params  *webhookpb.PublishLinkHitParams
		want    ExpectedRequests
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
			want: ExpectedRequests{
				numberOfRequests: 1,
				waitMilliseconds: 500,
				methods:          []string{"POST"},
				contentTypes:     []string{"application/json"},
				bodies: []string{
					setup.marshalJson(t,
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
								}}})}},
			wantErr: wantErr{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("test: %s", tt.name)
			setup.resetHttpRequests(t)
			got, err := client.PublishLinkHit(setup.ctx, tt.params)
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

func setupForPublishLinkHit(t *testing.T) *testDeps {
	setup := newSetup(t, tfx, listener)
	setup.startHttpServer(t)
	setup.startWorker(t)
	setup.adb.HaveInDatabase("webhook",
		"id, account_id, event, name, url, rate_limit, created_at, updated_at",
		[]interface{}{32767, "42", service.EventLinkHitStatus, "name1", setup.httpServer.URL, 2, "2020-01-12 22:41:42", "2020-01-12 22:41:42"})
	return setup
}
