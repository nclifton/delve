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

type ExpectedMMSStatusUpdateData = service.PublishStatusData
type ExpectedMMSStatusUpdateRequestBody struct {
	Event string                      `json:"event"`
	Data  ExpectedMMSStatusUpdateData `json:"data"`
}

func Test_PublishMMSStatusUpdate(t *testing.T) {
	log.Println("test PublishMMSStatusUpdate")

	setup := setupForPublishMMSStatusUpdate(t)
	defer setup.teardown(t)
	client := setup.getClient(t)
	timestampNow := timestamppb.Now()

	type wantErr struct {
		status *status.Status
		ok     bool
	}

	tests := []struct {
		name    string
		params  *webhookpb.PublishMMSStatusUpdateParams
		want    ExpectedRequests
		wantErr wantErr
	}{
		{
			name: "happy",
			params: &webhookpb.PublishMMSStatusUpdateParams{
				AccountId:         "42",
				MMSId:             "xxy",
				MessageRef:        "123",
				Recipient:         "35426378914",
				Sender:            "46354078643",
				Status:            "done",
				StatusDescription: "test is done",
				StatusUpdatedAt:   timestampNow,
			},
			want: ExpectedRequests{
				numberOfRequests: 1,
				waitMilliseconds: 500,
				methods:          []string{"POST"},
				contentTypes:     []string{"application/json"},
				bodies: []string{
					setup.marshalJson(t,
						ExpectedMMSStatusUpdateRequestBody{
							service.EventMMSStatus,
							ExpectedMMSStatusUpdateData{
								MMS_id:            "xxy",
								Message_ref:       "123",
								Recipient:         "35426378914",
								Sender:            "46354078643",
								Status:            "done",
								Status_updated_at: timestampNow.AsTime().Format(time.RFC3339),
							}}),
				},
			},
			wantErr: wantErr{},
		},
		{
			name: "unknown account id",
			params: &webhookpb.PublishMMSStatusUpdateParams{
				AccountId:         "43",
				MMSId:             "xxy",
				MessageRef:        "123",
				Recipient:         "35426378914",
				Sender:            "46354078643",
				Status:            "done",
				StatusDescription: "test is done",
				StatusUpdatedAt:   timestampNow,
			},
			want: ExpectedRequests{
				numberOfRequests: 0,
				waitMilliseconds: 500,
			},
			wantErr: wantErr{},
		},
		{
			name: "event not in webhooks for account",
			params: &webhookpb.PublishMMSStatusUpdateParams{
				AccountId:         "44",
				MMSId:             "xxy",
				MessageRef:        "123",
				Recipient:         "35426378914",
				Sender:            "46354078643",
				Status:            "done",
				StatusDescription: "test is done",
				StatusUpdatedAt:   timestampNow,
			},
			want: ExpectedRequests{
				numberOfRequests: 0,
				waitMilliseconds: 500,
			},
			wantErr: wantErr{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("test: %s", tt.name)
			setup.resetHttpRequests(t)
			got, err := client.PublishMMSStatusUpdate(setup.ctx, tt.params)
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

func setupForPublishMMSStatusUpdate(t *testing.T) *testDeps {
	setup := newSetup(t, tfx, listener)
	setup.startHttpServer(t)
	setup.startWorker(t)
	setup.adb.HaveInDatabase("webhook",
		"id, account_id, event, name, url, rate_limit, created_at, updated_at",
		[]interface{}{32767, "42", service.EventMMSStatus, "name1", setup.httpServer.URL, 2, "2020-01-12 22:41:42", "2020-01-12 22:41:42"})
	setup.adb.HaveInDatabase("webhook",
		"id, account_id, event, name, url, rate_limit, created_at, updated_at",
		[]interface{}{32768, "44", service.EventOptOutStatus, "name1", setup.httpServer.URL, 2, "2020-01-12 22:41:42", "2020-01-12 22:41:42"})
	return setup
}
