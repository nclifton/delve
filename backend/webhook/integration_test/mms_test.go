// +build integration

package test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/app/service"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
)

type ExpectedTimestamp struct {
	Nanos   int32 `json:"nanos"`
	Seconds int64 `json:"seconds"`
}
type ExpectedData struct {
	Account_id         string
	MMS_id             string
	Message_ref        string
	Recipient          string
	Sender             string
	Status             string
	Status_description string
	Status_updated_at  ExpectedTimestamp
}
type ExpectedRequestBody struct {
	Event string       `json:"event"`
	Data  ExpectedData `json:"data"`
}

func Test_PublishMMSStatusUpdate(t *testing.T) {
	log.Println("test PublishMMSStatusUpdate")

	setup := setupForPublishMMSStatusUpdate(t)
	defer setup.teardown(t)
	client := setup.getClient(t)

	type wantErr struct {
		status *status.Status
		ok     bool
	}
	timestampNow := timestamppb.Now()

	tests := []struct {
		name    string
		params  *webhookpb.PublishMMSStatusUpdateParams
		want    func(*webhookpb.NoReply, *webhookpb.PublishMMSStatusUpdateParams)
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
			want: func(reply *webhookpb.NoReply, params *webhookpb.PublishMMSStatusUpdateParams) {

				waitForRequest(setup, t)
				if len(setup.httpRequests) > 0 {
					req := setup.httpRequests[0]
					assert.Equal(t, req.Method, "POST", "request method")
					assert.Equal(t, req.Header.Get("Content-Type"), "application/json", "request has expected Content-Type")
					expectedBody, err := json.Marshal(
						ExpectedRequestBody{service.EventMMSStatus,
							ExpectedData{"42", "xxy", "123", "35426378914", "46354078643", "done", "test is done",
								ExpectedTimestamp{timestampNow.Nanos, timestampNow.Seconds}}})
					if err != nil {
						t.Fatalf("error: %+v", err)
					}
					assert.JSONEq(t, string(expectedBody), setup.httpRequestBodies[0], "request body json")
				}

			},
			wantErr: wantErr{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.PublishMMSStatusUpdate(setup.ctx, tt.params)
			if tt.wantErr.status != nil && err != nil {
				errStatus, ok := status.FromError(err)
				assert.Equal(t, ok, tt.wantErr.ok, "grpc ok")
				assert.EqualValues(t, tt.wantErr.status, errStatus, "grpc status")
			} else if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			tt.want(got, tt.params)
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
	return setup
}
