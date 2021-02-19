package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/burstsms/mtmo-tp/backend/sender/rpc/db"
)

func Test_senderImpl_validateCSVSender(t *testing.T) {

	ctx := context.TODO()

	type args struct {
		csvSender SenderCSV
	}
	type mock struct {
		method  string
		args    []interface{}
		returns []interface{}
	}
	type want struct {
		validSenders     []db.Sender
		validatedSenders []SenderCSV
	}

	tests := []struct {
		name string
		args args
		mock []mock
		want want
	}{
		{
			name: "blank address",
			args: args{
				csvSender: SenderCSV{
					AccountId:      "",
					Address:        "",
					Country:        "AU",
					Channels:       []string{"sms"},
					MMSProviderKey: "",
					Comment:        "",
					Status:         "",
					Error:          "",
				},
			},
			want: want{
				validSenders: []db.Sender{},
				validatedSenders: []SenderCSV{{
					AccountId:      "",
					Address:        "",
					Country:        "AU",
					Channels:       []string{"sms"},
					MMSProviderKey: "",
					Comment:        "",
					Status:         CSV_STATUS_SKIPPED,
					Error:          "Address: required",
				}},
			},
		},
		{
			name: "address not unique",
			args: args{
				csvSender: SenderCSV{
					AccountId:      "",
					Address:        "RHINO",
					Country:        "AU",
					Channels:       []string{"sms"},
					MMSProviderKey: "",
					Comment:        "",
					Status:         "",
					Error:          "",
				},
			},
			mock: []mock{
				{
					method: "FindSendersByAddress",
					args:   []interface{}{ctx, "RHINO"},
					returns: []interface{}{[]db.Sender{{
						Address: "RHINO",
					}}, nil},
				},
			},
			want: want{
				validSenders: []db.Sender{},
				validatedSenders: []SenderCSV{{
					AccountId:      "",
					Address:        "RHINO",
					Country:        "AU",
					Channels:       []string{"sms"},
					MMSProviderKey: "",
					Comment:        "",
					Status:         CSV_STATUS_SKIPPED,
					Error:          "Address: exists",
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := new(db.MockDB)
			for _, mock := range tt.mock {
				mockDb.On(mock.method, mock.args...).Return(mock.returns...)
			}
			s := &senderImpl{
				db: mockDb,
			}
			validSenders, validatedSenders := s.validateCSVSenders(ctx, []SenderCSV{tt.args.csvSender})
			assert.Equal(t, tt.want.validatedSenders, validatedSenders, "validated senders")
			assert.Equal(t, tt.want.validSenders, validSenders, "valid senders")
		})
	}
}
