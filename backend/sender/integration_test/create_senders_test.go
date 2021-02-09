// +build integration

package test

import (
	"testing"

	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"

	"github.com/burstsms/mtmo-tp/backend/lib/assertdb"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
)

func Test_CreateSenders(t *testing.T) {
	s := setupForCreate(t)
	defer s.teardown(t)
	client := s.getClient(t)

	type want struct {
		reply *senderpb.CreateSendersReply
	}

	type wantErr struct {
		status *status.Status
		ok     bool
	}

	tests := []struct {
		name    string
		params  *senderpb.CreateSendersParams
		want    want
		wantErr wantErr
	}{
		{
			name: "happy create one",
			params: &senderpb.CreateSendersParams{
				Senders: []*senderpb.NewSender{
					{
						AccountId:      "",
						Address:        "LION",
						MMSProviderKey: "fake",
						Channels:       []string{"mms", "sms"},
						Country:        "AU",
						Comment:        "roars",
					},
				},
			},
			want: want{
				reply: &senderpb.CreateSendersReply{}},
			wantErr: wantErr{},
		}, {
			name: "happy create two",
			params: &senderpb.CreateSendersParams{
				Senders: []*senderpb.NewSender{
					{
						AccountId:      "",
						Address:        "TIGER",
						MMSProviderKey: "optus",
						Channels:       []string{"mms"},
						Country:        "AU",
						Comment:        "roars too",
					},{
						AccountId:      "",
						Address:        "PANTHER",
						MMSProviderKey: "mgage",
						Channels:       []string{"sms"},
						Country:        "US",
						Comment:        "is silent",
					},
				},
			},
			want: want{
				reply: &senderpb.CreateSendersReply{}},
			wantErr: wantErr{},
		},
		{
			name: "empty",
			params: &senderpb.CreateSendersParams{
				Senders: []*senderpb.NewSender{},
			},
			want: want{
				reply: &senderpb.CreateSendersReply{}},
			wantErr: wantErr{},
		},
		{
			name: "unvalidated sender - blank values",
			params: &senderpb.CreateSendersParams{
				Senders: []*senderpb.NewSender{{
					AccountId:      "",
					Address:        "",
					MMSProviderKey: "",
					Channels:       []string{},
					Country:        "",
					Comment:        "",
				}},
			},
			want: want{
				reply: &senderpb.CreateSendersReply{}},
			wantErr: wantErr{},
		},
		{
			name: "unvalidated sender - nil values",
			params: &senderpb.CreateSendersParams{
				Senders: []*senderpb.NewSender{{}},
			},
			want: want{
				reply: &senderpb.CreateSendersReply{}},
			wantErr: wantErr{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.CreateSenders(s.ctx, tt.params)
			if tt.wantErr.status != nil && err != nil {
				errStatus, ok := status.FromError(err)
				assert.Equal(t, ok, tt.wantErr.ok, "grpc ok")
				assert.EqualValues(t, tt.wantErr.status, errStatus, "grpc status")
			} else if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			} else {
				assert.Equal(t, len(tt.want.reply.Senders), len(got.Senders), "number of senders in reply")
				for idx, sender := range tt.want.reply.Senders {
					assert.NotEmptyf(t, sender.Id, "sender %d Id", idx)
					assert.Equalf(t, sender.AccountId, got.Senders[idx].AccountId, "sender %d AccountId", idx)
					assert.Equalf(t, sender.Address, got.Senders[idx].Address, "sender %d Address", idx)
					assert.Equalf(t, sender.Channels, got.Senders[idx].Channels, "sender %d Channels", idx)
					assert.Equalf(t, sender.MMSProviderKey, got.Senders[idx].MMSProviderKey, "sender %d MMSProviderKey", idx)
					assert.Equalf(t, sender.Country, got.Senders[idx].Country, "sender %d Country", idx)
					assert.Equalf(t, sender.Comment, got.Senders[idx].Comment, "sender %d Channels", idx)
					assert.NotEmptyf(t, sender.CreatedAt, "sender %d CreatedAt", idx)
					assert.NotEmptyf(t, sender.UpdatedAt, "sender %d UpdatedAt", idx)

					s.SeeInDatabase("sender", assertdb.Criteria{
						"account_id":       sender.AccountId,
						"address":          sender.Address,
						"mms_provider_key": sender.MMSProviderKey,
						"channels":         sender.Channels,
						"country":          sender.Country,
						"comment":          sender.Comment,
					})

				}
			}
		})
	}
}

func setupForCreate(t *testing.T) *testDeps {
	s := newSetup(t, tfx)

	return s
}
