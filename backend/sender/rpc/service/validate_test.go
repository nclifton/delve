package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/burstsms/mtmo-tp/backend/sender/rpc/db"
)

func Test_senderImpl_validateCSVSender(t *testing.T) {

	ctx := context.TODO()
	senderEnums := db.SenderEnums{
		"provider_key": []string{"fake", "optus", "mgage"},
		"channel":      []string{"mms", "sms"},
	}

	type args struct {
		csvSenders []SenderCSV
	}
	type mock struct {
		method  string
		args    []interface{}
		returns []interface{}
		times   int
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
				[]SenderCSV{{"", "", "AU", []string{"sms"}, "", "", "", ""}},
			},
			mock: []mock{
				{"GetSenderEnums", []interface{}{ctx}, []interface{}{senderEnums, nil}, 1},
			},
			want: want{
				[]db.Sender{},
				[]SenderCSV{
					{"", "", "AU", []string{"sms"}, "", "", CSV_STATUS_SKIPPED, "Address: required"},
				},
			},
		},
		{
			name: "address not new",
			args: args{
				[]SenderCSV{{"", "RHINO", "AU", []string{"sms"}, "", "", "", ""}},
			},
			mock: []mock{
				{"GetSenderEnums", []interface{}{ctx}, []interface{}{senderEnums, nil}, 1},
				{"SenderAddressExists", []interface{}{ctx, "RHINO"}, []interface{}{true, nil}, 1},
			},
			want: want{
				[]db.Sender{},
				[]SenderCSV{
					{"", "RHINO", "AU", []string{"sms"}, "", "", CSV_STATUS_SKIPPED, "Address: is not new"},
				},
			},
		}, {
			name: "country required",
			args: args{
				[]SenderCSV{{"", "RHINO", "", []string{"sms"}, "", "", "", ""}},
			},
			mock: []mock{
				{"GetSenderEnums", []interface{}{ctx}, []interface{}{senderEnums, nil}, 1},
				{"SenderAddressExists", []interface{}{ctx, "RHINO"}, []interface{}{false, nil}, 1},
			},
			want: want{
				[]db.Sender{},
				[]SenderCSV{
					{"", "RHINO", "", []string{"sms"}, "", "", CSV_STATUS_SKIPPED, "Country: required"},
				},
			},
		}, {
			name: "mms provider key is not one of enum",
			args: args{
				[]SenderCSV{{"", "RHINO", "AU", []string{"sms"}, "bad", "", "", ""}},
			},
			mock: []mock{
				{"GetSenderEnums", []interface{}{ctx}, []interface{}{senderEnums, nil}, 2},
				{"SenderAddressExists", []interface{}{ctx, "RHINO"}, []interface{}{false, nil}, 1},
			},
			want: want{
				[]db.Sender{},
				[]SenderCSV{
					{"", "RHINO", "AU", []string{"sms"}, "bad", "", CSV_STATUS_SKIPPED, "MMSProviderKey: bad did not match any of fake,optus,mgage"},
				},
			},
		}, {
			name: "channels required",
			args: args{
				[]SenderCSV{{"", "RHINO", "AU", []string{}, "", "", "", ""}},
			},
			mock: []mock{
				{"SenderAddressExists", []interface{}{ctx, "RHINO"}, []interface{}{false, nil}, 1},
			},
			want: want{
				[]db.Sender{},
				[]SenderCSV{
					{"", "RHINO", "AU", []string{}, "", "", CSV_STATUS_SKIPPED, "Channels: required"},
				},
			},
		}, {
			name: "channels is not one of",
			args: args{
				[]SenderCSV{{"", "RHINO", "AU", []string{"bad"}, "", "", "", ""}},
			},
			mock: []mock{
				{"GetSenderEnums", []interface{}{ctx}, []interface{}{senderEnums, nil}, 1},
				{"SenderAddressExists", []interface{}{ctx, "RHINO"}, []interface{}{false, nil}, 1},
			},
			want: want{
				[]db.Sender{},
				[]SenderCSV{
					{"", "RHINO", "AU", []string{"bad"}, "", "", CSV_STATUS_SKIPPED, "Channels: bad did not match any of mms,sms"},
				},
			},
		}, {
			name: "country is not one of",
			args: args{
				[]SenderCSV{{"", "RHINO", "DE", []string{"sms"}, "", "", "", ""}},
			},
			mock: []mock{
				{"GetSenderEnums", []interface{}{ctx}, []interface{}{senderEnums, nil}, 1},
				{"SenderAddressExists", []interface{}{ctx, "RHINO"}, []interface{}{false, nil}, 1},
			},
			want: want{
				[]db.Sender{},
				[]SenderCSV{
					{"", "RHINO", "DE", []string{"sms"}, "", "", CSV_STATUS_SKIPPED, "Country: DE did not match any of au,AU,us,US"},
				},
			},
		}, {
			name: "mms provider key required if channels contains mms - error",
			args: args{
				[]SenderCSV{{"", "RHINO", "AU", []string{"mms"}, "", "", "", ""}},
			},
			mock: []mock{
				{"GetSenderEnums", []interface{}{ctx}, []interface{}{senderEnums, nil}, 1},
				{"SenderAddressExists", []interface{}{ctx, "RHINO"}, []interface{}{false, nil}, 1},
			},
			want: want{
				[]db.Sender{},
				[]SenderCSV{
					{"", "RHINO", "AU", []string{"mms"}, "", "", CSV_STATUS_SKIPPED, "MMSProviderKey: required if Channels contains mms"},
				},
			},
		}, {
			name: "mms provider key required if channels contains mms - ok",
			args: args{
				[]SenderCSV{{"", "RHINO", "AU", []string{"mms"}, "fake", "", "", ""}},
			},
			mock: []mock{
				{"GetSenderEnums", []interface{}{ctx}, []interface{}{senderEnums, nil}, 2},
				{"SenderAddressExists", []interface{}{ctx, "RHINO"}, []interface{}{false, nil}, 1},
			},
			want: want{
				[]db.Sender{{
					Address:        "RHINO",
					MMSProviderKey: "fake",
					Channels:       []string{"mms"},
					Country:        "AU",
					Comment:        "",
				}},
				[]SenderCSV{
					{"", "RHINO", "AU", []string{"mms"}, "fake", "", CSV_STATUS_OK, ""},
				},
			},
		}, {
			name: "all address occupancies skipped if occurs more than once",
			args: args{
				[]SenderCSV{
					{"", "RHINO", "AU", []string{"mms"}, "fake", "", "", ""},
					{"", "RHINO", "AU", []string{"mms"}, "fake", "", "", ""},
				},
			},
			mock: []mock{
				{"GetSenderEnums", []interface{}{ctx}, []interface{}{senderEnums, nil}, 4},
				{"SenderAddressExists", []interface{}{ctx, "RHINO"}, []interface{}{false, nil}, 2},
			},
			want: want{
				[]db.Sender{},
				[]SenderCSV{
					{"", "RHINO", "AU", []string{"mms"}, "fake", "", CSV_STATUS_SKIPPED, "Address: multiple occurrence in upload"},
					{"", "RHINO", "AU", []string{"mms"}, "fake", "", CSV_STATUS_SKIPPED, "Address: multiple occurrence in upload"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := new(db.MockDB)
			for _, mock := range tt.mock {
				mockDb.On(mock.method, mock.args...).Return(mock.returns...).Times(mock.times)
			}
			s := &senderImpl{
				db: mockDb,
			}
			validSenders, validatedSenders := s.validateCSVSenders(ctx, tt.args.csvSenders)
			assert.Equal(t, tt.want.validatedSenders, validatedSenders, "validated senders")
			assert.Equal(t, tt.want.validSenders, validSenders, "valid senders")
			mockDb.AssertExpectations(t)
		})
	}
}
