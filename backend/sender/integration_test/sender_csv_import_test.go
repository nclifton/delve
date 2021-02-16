// +build integration

package test

import (
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"
	"github.com/vincent-petithory/dataurl"

	"github.com/burstsms/mtmo-tp/backend/lib/assertdb"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
)

func Test_CreateSendersFromCSVDataURL(t *testing.T) {
	s := setupForCreate(t)
	defer s.teardown(t)
	client := s.getClient(t)

	type want struct {
		reply      *senderpb.CreateSendersFromCSVDataURLReply
		dbCriteria []assertdb.Criteria
	}

	type wantErr struct {
		status *status.Status
		ok     bool
	}

	type args struct {
		csv string
	}

	tests := []struct {
		name    string
		args    args
		want    want
		wantErr wantErr
	}{
		{
			name: "happy create one",
			args: args{
				csv: `account_id,address,country,channels,mms_provider_key,comment
					,LION,AU,"[""mms"",""sms""]",fake,"roars"`,
			},
			want: want{
				reply: &senderpb.CreateSendersFromCSVDataURLReply{
					Senders: []*senderpb.Sender{{
						AccountId:      "",
						Address:        "LION",
						MMSProviderKey: "fake",
						Channels:       []string{"mms", "sms"},
						Country:        "AU",
						Comment:        "roars",
					}},
				},
				dbCriteria: []map[string]interface{}{
					{
						"account_id":       nil,
						"address":          "LION",
						"mms_provider_key": "fake",
						"channels":         []string{"mms", "sms"},
						"country":          "AU",
						"comment":          "roars",
					},
				},
			},

			wantErr: wantErr{},
		},
		 {
			name: "happy create two",
			args: args{
				csv: `account_id,address,country,channels,mms_provider_key,comment
					,TIGER,AU,"[""mms""]",optus,"roars too"
					,PANTHER,US,"[""sms""]",mgage,"is silent"`,
			},
			want: want{
				reply: &senderpb.CreateSendersFromCSVDataURLReply{
					Senders: []*senderpb.Sender{{
						AccountId:      "",
						Address:        "TIGER",
						MMSProviderKey: "optus",
						Channels:       []string{"mms"},
						Country:        "AU",
						Comment:        "roars too",
					}, {
						AccountId:      "",
						Address:        "PANTHER",
						MMSProviderKey: "mgage",
						Channels:       []string{"sms"},
						Country:        "US",
						Comment:        "is silent",
					}},
				},
				dbCriteria: []map[string]interface{}{
					{
						"account_id":       nil,
						"address":          "TIGER",
						"mms_provider_key": "optus",
						"channels":         []string{"mms"},
						"country":          "AU",
						"comment":          "roars too",
					}, {
						"account_id":       nil,
						"address":          "PANTHER",
						"mms_provider_key": "mgage",
						"channels":         []string{"sms"},
						"country":          "US",
						"comment":          "is silent",
					},
				},
			},
			wantErr: wantErr{},
		},
		{
			name: "empty",
			args: args{
				csv: `account_id,address,country,channels,mms_provider_key,comment`,
			},
			want: want{
				reply:      &senderpb.CreateSendersFromCSVDataURLReply{},
				dbCriteria: []map[string]interface{}{},
			},
			wantErr: wantErr{},
		},
		{
			name: "unvalidated sender - blank values",
			args: args{
				csv: `account_id,address,country,channels,mms_provider_key,comment
					,,,"[]",,`,
			},
			want: want{
				reply: &senderpb.CreateSendersFromCSVDataURLReply{},
			},
			wantErr: wantErr{
				status: status.New(codes.Unknown, `ERROR: null value in column "address" violates not-null constraint (SQLSTATE 23502)`),
				ok:     true, // the grpc service did respond to the call
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			du, err := dataurl.New([]byte(tt.args.csv), "text/csv").MarshalText()
			if err != nil {
				t.Fatal(err)
			}
			params := &senderpb.CreateSendersFromCSVDataURLParams{
				CSV: du,
			}
			got, err := client.CreateSendersFromCSVDataURL(s.ctx, params)
			if tt.wantErr.status != nil {
				if err != nil {
					errStatus, ok := status.FromError(err)
					assert.Equal(t, ok, tt.wantErr.ok, "grpc ok")
					assert.EqualValues(t, tt.wantErr.status, errStatus, "grpc status")
				} else {
					t.Fatal("did not get expected error")
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			} else {
				assert.Equal(t, len(tt.want.reply.Senders), len(got.Senders), "number of senders in reply")
				for idx, sender := range tt.want.reply.Senders {
					assert.NotEmptyf(t, got.Senders[idx].Id, "sender %d Id", idx)
					assert.Equalf(t, sender.AccountId, got.Senders[idx].AccountId, "sender %d AccountId", idx)
					assert.Equalf(t, sender.Address, got.Senders[idx].Address, "sender %d Address", idx)
					assert.Equalf(t, sender.Channels, got.Senders[idx].Channels, "sender %d Channels", idx)
					assert.Equalf(t, sender.MMSProviderKey, got.Senders[idx].MMSProviderKey, "sender %d MMSProviderKey", idx)
					assert.Equalf(t, sender.Country, got.Senders[idx].Country, "sender %d Country", idx)
					assert.Equalf(t, sender.Comment, got.Senders[idx].Comment, "sender %d Channels", idx)
					assert.NotEmptyf(t, got.Senders[idx].CreatedAt, "sender %d CreatedAt", idx)
					assert.NotEmptyf(t, got.Senders[idx].UpdatedAt, "sender %d UpdatedAt", idx)

				}
				for _, criteria := range tt.want.dbCriteria {
					s.SeeInDatabase("sender", criteria)
				}

			}
		})
	}
}

func setupForCreate(t *testing.T) *testDeps {
	s := newSetup(t, tfx)

	return s
}
