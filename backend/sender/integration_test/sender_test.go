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
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/app/run"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
)

var tfx *fixtures.TestFixtures

func TestMain(m *testing.M) {
	tfx = fixtures.New()
	tfx.SetupPostgres("sender")
	tfx.GRPCStart(run.Server)
	code := m.Run()
	defer os.Exit(code)
	defer tfx.Teardown()
}

func Test_FindByAddress(t *testing.T) {
	s := setupForFind(t)
	defer s.teardown(t)
	client := s.getClient(t)

	type want struct {
		reply *senderpb.FindByAddressReply
	}

	type wantErr struct {
		status *status.Status
		ok     bool
	}

	tests := []struct {
		name    string
		params  *senderpb.FindByAddressParams
		want    want
		wantErr wantErr
	}{
		{
			name: "happy",
			params: &senderpb.FindByAddressParams{
				AccountId: s.uuids[1],
				Address:   "FISH",
			},
			want: want{
				reply: &senderpb.FindByAddressReply{
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
			name: "not found Address",
			params: &senderpb.FindByAddressParams{
				AccountId: s.uuids[1],
				Address:   "CHIPS",
			},
			want: want{reply: nil},
			wantErr: wantErr{
				status: status.New(codes.NotFound, "sender not found"),
				ok:     true, // the grpc service did respond
			},
		},
		{
			name: "not found Account",
			params: &senderpb.FindByAddressParams{
				AccountId: s.uuids[2],
				Address:   "FISH",
			},
			want: want{reply: nil},
			wantErr: wantErr{
				status: status.New(codes.NotFound, "sender not found"),
				ok:     true, // the grpc service did respond
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.FindByAddress(s.ctx, tt.params)
			if tt.wantErr.status != nil && err != nil {
				errStatus, ok := status.FromError(err)
				assert.Equal(t, ok, tt.wantErr.ok, "grpc ok")
				assert.EqualValues(t, tt.wantErr.status, errStatus, "grpc status")
			} else if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			} else {
				equal := assert.ObjectsAreEqual(tt.want.reply.Sender, got.Sender)
				assert.True(t, equal, fmt.Sprintf("reply \n\twant: \n\t%+v\n\tgot: \n\t%+v\n", tt.want.reply.Sender, got.Sender))
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

	return s
}
