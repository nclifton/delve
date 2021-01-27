package biz_test

import (
	"testing"
	"time"

	"github.com/burstsms/mtmo-tp/backend/lib/errorlib"
	"github.com/burstsms/mtmo-tp/backend/mms/biz"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_IsValidSender(t *testing.T) {
	type args struct {
		sender  *senderpb.Sender
		address string
		country string
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "not found",
			args: args{
				sender:  nil,
				address: "BOOB",
				country: "au",
			},
			want: errorlib.ErrInvalidSenderNotFound,
		},
		{
			name: "invalid address",
			args: args{
				sender: &senderpb.Sender{
					Id:             "12345",
					AccountId:      "23456",
					Address:        "OBBO",
					MMSProviderKey: "optus",
					Channels:       []string{"mms", "sms"},
					Country:        "au",
					Comment:        "blah",
					CreatedAt:      timestamppb.New(time.Now()),
					UpdatedAt:      timestamppb.New(time.Now()),
				},
				address: "BOOB",
				country: "au",
			},
			want: errorlib.ErrInvalidSenderAddress,
		},
		{
			name: "invalid channel",
			args: args{
				sender: &senderpb.Sender{
					Id:             "12345",
					AccountId:      "23456",
					Address:        "BOOB",
					MMSProviderKey: "optus",
					Channels:       []string{"sms"},
					Country:        "au",
					Comment:        "blah",
					CreatedAt:      timestamppb.New(time.Now()),
					UpdatedAt:      timestamppb.New(time.Now()),
				},
				address: "BOOB",
				country: "au",
			},
			want: errorlib.ErrInvalidSenderChannel,
		},
		{
			name: "happy",
			args: args{
				sender: &senderpb.Sender{
					Id:             "12345",
					AccountId:      "23456",
					Address:        "BOOB",
					MMSProviderKey: "optus",
					Channels:       []string{"mms", "sms"},
					Country:        "au",
					Comment:        "blah",
					CreatedAt:      timestamppb.New(time.Now()),
					UpdatedAt:      timestamppb.New(time.Now()),
				},
				address: "BOOB",
				country: "au",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := biz.IsValidSender(tt.args.sender, tt.args.address, tt.args.country)
			if err != tt.want {
				t.Errorf("expected: %s\n got: %s\n", tt.want, err)
			}
		})
	}
}
