// +build integration

package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burstsms/mtmo-tp/backend/lib/assertdb"
	"github.com/burstsms/mtmo-tp/backend/lib/fixtures"
	"github.com/burstsms/mtmo-tp/backend/lib/rpcbuilder"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/builder"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
)

var tfx *fixtures.TestFixtures

func TestMain(m *testing.M) {
	tfx = fixtures.New()
	tfx.SetupPostgres("sender")
	tfx.GRPCStart(senderRPCService())
	code := m.Run()
	defer os.Exit(code)
	defer tfx.Teardown()
}

func senderRPCService() rpcbuilder.Service {
	return builder.NewBuilderFromEnv()
}

func Test_FindSenderByAddressAndAccountID(t *testing.T) {
	s := setupForFind(t)
	defer s.teardown(t)
	client := s.getClient(t)

	type want struct {
		reply *senderpb.FindSenderByAddressAndAccountIDReply
	}

	type wantErr struct {
		status *status.Status
		ok     bool
	}

	tests := []struct {
		name    string
		params  *senderpb.FindSenderByAddressAndAccountIDParams
		want    want
		wantErr wantErr
	}{
		{
			name: "happy",
			params: &senderpb.FindSenderByAddressAndAccountIDParams{
				AccountId: s.uuids[1],
				Address:   "FISH",
			},
			want: want{
				reply: &senderpb.FindSenderByAddressAndAccountIDReply{
					Sender: &senderpb.Sender{
						Id:             s.uuids[0],
						AccountId:      s.uuids[1],
						Address:        "FISH",
						MMSProviderKey: "optus",
						Channels:       []string{"mms", "sms"},
						Country:        "AU",
						Comment:        "Slartibartfast",
						CreatedAt:      timestamppb.New(s.dates[0]),
						UpdatedAt:      timestamppb.New(s.dates[0]),
					}}},
			wantErr: wantErr{},
		},
		{
			name: "no comment",
			params: &senderpb.FindSenderByAddressAndAccountIDParams{
				AccountId: s.uuids[6],
				Address:   "NOCOMMENT",
			},
			want: want{&senderpb.FindSenderByAddressAndAccountIDReply{
				Sender: &senderpb.Sender{
					Id:             s.uuids[5],
					AccountId:      s.uuids[6],
					Address:        "NOCOMMENT",
					MMSProviderKey: "optus",
					Channels:       []string{"mms", "sms"},
					Country:        "AU",
					Comment:        "",
					CreatedAt:      timestamppb.New(s.dates[5]),
					UpdatedAt:      timestamppb.New(s.dates[6]),
				}}},
			wantErr: wantErr{},
		},

		{
			name: "not found sender Address: MICE",
			params: &senderpb.FindSenderByAddressAndAccountIDParams{
				AccountId: s.uuids[1],
				Address:   "MICE",
			},
			want: want{
				reply: &senderpb.FindSenderByAddressAndAccountIDReply{
					Sender: nil,
				},
			},
			wantErr: wantErr{
				status: status.New(codes.NotFound, "sender not found"),
				ok:     true, // the grpc service did respond
			},
		},
		{
			name: fmt.Sprintf("not found sender Account: %s", s.uuids[2]),
			params: &senderpb.FindSenderByAddressAndAccountIDParams{
				AccountId: s.uuids[2],
				Address:   "FISH",
			},
			want: want{
				reply: &senderpb.FindSenderByAddressAndAccountIDReply{
					Sender: nil,
				}},
			wantErr: wantErr{
				status: status.New(codes.NotFound, "sender not found"),
				ok:     true, // the grpc service did respond
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.FindSenderByAddressAndAccountID(s.ctx, tt.params)
			if tt.wantErr.status != nil && err != nil {
				errStatus, ok := status.FromError(err)
				assert.Equal(t, ok, tt.wantErr.ok, "grpc ok")
				assert.EqualValues(t, tt.wantErr.status, errStatus, "grpc status")
			} else if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			} else {
				equal := assert.ObjectsAreEqual(tt.want.reply.Sender, got.Sender)
				assert.True(t, equal, fmt.Sprintf("reply sender \n\twant: \n\t%+v\n\tgot: \n\t%+v\n", tt.want.reply.Sender, got.Sender))
			}

		})
	}
}

func Test_FindSendersByAddress(t *testing.T) {
	s := setupForFind(t)
	defer s.teardown(t)
	client := s.getClient(t)

	type wantErr struct {
		status *status.Status
		ok     bool
	}

	tests := []struct {
		name    string
		params  *senderpb.FindSendersByAddressParams
		want    *senderpb.FindSendersByAddressReply
		wantErr wantErr
	}{
		{
			name: "happy",
			params: &senderpb.FindSendersByAddressParams{
				Address: "FISH",
			},
			want: &senderpb.FindSendersByAddressReply{
				Senders: []*senderpb.Sender{
					{
						Id:             s.uuids[0],
						AccountId:      s.uuids[1],
						Address:        "FISH",
						MMSProviderKey: "optus",
						Channels:       []string{"mms", "sms"},
						Country:        "AU",
						Comment:        "Slartibartfast",
						CreatedAt:      timestamppb.New(s.dates[0]),
						UpdatedAt:      timestamppb.New(s.dates[0]),
					},
				}},
			wantErr: wantErr{},
		},
		{
			name: "address not found",
			params: &senderpb.FindSendersByAddressParams{
				Address: "NOTFOUND",
			},
			want:    &senderpb.FindSendersByAddressReply{Senders: []*senderpb.Sender{}},
			wantErr: wantErr{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.FindSendersByAddress(s.ctx, tt.params)
			if tt.wantErr.status != nil && err != nil {
				errStatus, ok := status.FromError(err)
				assert.Equal(t, ok, tt.wantErr.ok, "grpc ok")
				assert.EqualValues(t, tt.wantErr.status, errStatus, "grpc status")
			} else if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			} else {
				assert.Equal(t, len(tt.want.Senders), len(got.Senders), "number of senders in reply")
				for idx, sender := range tt.want.Senders {
					equal := assert.ObjectsAreEqual(sender, got.Senders[idx])
					assert.True(t, equal, fmt.Sprintf("reply Senders[%d] \n\twant: \n\t%+v\n\tgot: \n\t%+v\n", idx, sender, got.Senders[idx]))
				}
			}
		})
	}
}

func Test_FindSendersByAccountId(t *testing.T) {
	s := setupForFind(t)
	defer s.teardown(t)
	client := s.getClient(t)

	type want struct {
		reply *senderpb.FindSendersByAccountIdReply
	}

	type wantErr struct {
		status *status.Status
		ok     bool
	}

	tests := []struct {
		name    string
		params  *senderpb.FindSendersByAccountIdParams
		want    want
		wantErr wantErr
	}{
		{
			name: "happy",
			params: &senderpb.FindSendersByAccountIdParams{
				AccountId: s.uuids[1],
			},
			want: want{
				reply: &senderpb.FindSendersByAccountIdReply{
					Senders: []*senderpb.Sender{
						{
							Id:             s.uuids[0],
							AccountId:      s.uuids[1],
							Address:        "FISH",
							MMSProviderKey: "optus",
							Channels:       []string{"mms", "sms"},
							Country:        "AU",
							Comment:        "Slartibartfast",
							CreatedAt:      timestamppb.New(s.dates[0]),
							UpdatedAt:      timestamppb.New(s.dates[0]),
						},
						{
							Id:             s.uuids[2],
							AccountId:      s.uuids[1],
							Address:        "CHIPS",
							MMSProviderKey: "optus",
							Channels:       []string{"mms", "sms"},
							Country:        "AU",
							Comment:        "Arthur Dent",
							CreatedAt:      timestamppb.New(s.dates[2]),
							UpdatedAt:      timestamppb.New(s.dates[1]),
						},
					},
				}},
			wantErr: wantErr{},
		},
		{
			name: "happy no comment",
			params: &senderpb.FindSendersByAccountIdParams{
				AccountId: s.uuids[6],
			},
			want: want{
				reply: &senderpb.FindSendersByAccountIdReply{
					Senders: []*senderpb.Sender{
						{
							Id:             s.uuids[5],
							AccountId:      s.uuids[6],
							Address:        "NOCOMMENT",
							MMSProviderKey: "optus",
							Channels:       []string{"mms", "sms"},
							Country:        "AU",
							Comment:        "",
							CreatedAt:      timestamppb.New(s.dates[5]),
							UpdatedAt:      timestamppb.New(s.dates[6]),
						},
					},
				}},
			wantErr: wantErr{},
		},
		{
			name: "account id not found",
			params: &senderpb.FindSendersByAccountIdParams{
				AccountId: s.uuids[5],
			},
			want: want{
				reply: &senderpb.FindSendersByAccountIdReply{Senders: []*senderpb.Sender{}}},
			wantErr: wantErr{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.FindSendersByAccountId(s.ctx, tt.params)
			if tt.wantErr.status != nil && err != nil {
				errStatus, ok := status.FromError(err)
				assert.Equal(t, ok, tt.wantErr.ok, "grpc ok")
				assert.EqualValues(t, tt.wantErr.status, errStatus, "grpc status")
			} else if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			} else {
				assert.Equal(t, len(tt.want.reply.Senders), len(got.Senders), "number of senders in reply")
				for idx, sender := range tt.want.reply.Senders {
					equal := assert.ObjectsAreEqual(sender, got.Senders[idx])
					assert.True(t, equal, fmt.Sprintf("reply Senders[%d] \n\twant: \n\t%+v\n\tgot: \n\t%+v\n", idx, sender, got.Senders[idx]))
				}
			}
		})
	}
}

func setupForFind(t *testing.T) *testDeps {
	s := newSetup(t, tfx)

	s.HaveInDatabase("sender", assertdb.Row{
		"id":               s.uuids[0],
		"account_id":       s.uuids[1],
		"address":          "FISH",
		"mms_provider_key": "optus",
		"channels":         []string{"mms", "sms"},
		"country":          "AU",
		"comment":          "Slartibartfast",
		"created_at":       s.dates[0],
		"updated_at":       s.dates[0]})

	s.HaveInDatabase("sender", assertdb.Row{
		"id":               s.uuids[2],
		"account_id":       s.uuids[1],
		"address":          "CHIPS",
		"mms_provider_key": "optus",
		"channels":         []string{"mms", "sms"},
		"country":          "AU",
		"comment":          "Arthur Dent",
		"created_at":       s.dates[2],
		"updated_at":       s.dates[1]})

	s.HaveInDatabase("sender", assertdb.Row{
		"id":               s.uuids[3],
		"account_id":       s.uuids[4],
		"address":          "MICE",
		"mms_provider_key": "mgage",
		"channels":         []string{"mms"},
		"country":          "AU",
		"comment":          "Arthur Dent",
		"created_at":       s.dates[4],
		"updated_at":       s.dates[3]})

	s.HaveInDatabase("sender", assertdb.Row{
		"id":               s.uuids[5],
		"account_id":       s.uuids[6],
		"address":          "NOCOMMENT",
		"mms_provider_key": "optus",
		"channels":         []string{"mms", "sms"},
		"country":          "AU",
		"created_at":       s.dates[5],
		"updated_at":       s.dates[6]})

	return s
}
