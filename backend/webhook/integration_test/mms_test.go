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

	i := setupForPublishMMSStatusUpdate(t)
	defer i.teardown(t)
	client := i.getClient(t)
	timestampNow := timestamppb.Now()

	type want struct {
		reply *webhookpb.NoReply
		requests ExpectedRequests
	}

	type wantErr struct {
		status *status.Status
		ok     bool
	}

	tests := []struct {
		name   string
		params *webhookpb.PublishMMSStatusUpdateParams
		want want
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
			want: want{
				reply: &webhookpb.NoReply{},
				requests: ExpectedRequests{
						NumberOfRequests: 1,
						WaitMilliseconds: 500,
						Methods:          []string{"POST"},
						ContentTypes:     []string{"application/json"},
						Bodies: []string{
							jsonString(t,
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
						}}},
			wantErr: wantErr{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("test: %s", tt.name)
			i.ResetHttpRequests()
			got, err := client.PublishMMSStatusUpdate(i.ctx, tt.params)
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

func setupForPublishMMSStatusUpdate(t *testing.T) *testDeps {
	i := newSetup(t, tfx)
	i.HaveInDatabase("webhook",
		"id, account_id, event, name, url, rate_limit, created_at, updated_at",
		[]interface{}{32767, "42", service.EventMMSStatus, "name1", i.webhookURL, 2, "2020-01-12 22:41:42", "2020-01-12 22:41:42"})
	i.HaveInDatabase("webhook",
		"id, account_id, event, name, url, rate_limit, created_at, updated_at",
		[]interface{}{32768, "44", service.EventOptOutStatus, "name1", i.webhookURL, 2, "2020-01-12 22:41:42", "2020-01-12 22:41:42"})
	return i
}
