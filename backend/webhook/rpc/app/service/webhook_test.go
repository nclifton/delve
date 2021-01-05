package service

import (
	"reflect"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/app/db"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
)

func Test_dbWebhookToWebhook(t *testing.T) {
	createdAt := time.Now().Add(-time.Hour)
	updatedAt := time.Now()
	type args struct {
		w db.Webhook
	}
	tests := []struct {
		name string
		args args
		want *webhookpb.Webhook
	}{
		// TODO: Add test cases.

		{
			name: "mapping",
			args: args{
				w: db.Webhook{
					ID:        1,
					AccountID: "1",
					Event:     "event",
					Name:      "name",
					URL:       "url",
					RateLimit: 1,
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
			},
			want: &webhookpb.Webhook{
				Id:        1,
				AccountId: "1",
				Event:     "event",
				Name:      "name",
				URL:       "url",
				RateLimit: 1,
				CreatedAt: timestamppb.New(createdAt),
				UpdatedAt: timestamppb.New(updatedAt),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dbWebhookToWebhook(tt.args.w); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dbWebhookToWebhook() = %v, want %v", got, tt.want)
			}
		})
	}
}
