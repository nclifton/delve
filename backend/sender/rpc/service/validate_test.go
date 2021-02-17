package service

import (
	"context"
	"reflect"
	"testing"

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

	tests := []struct {
		name       string
		args       args
		mock       []mock
		wantValid  bool
		wantResult SenderCSV
		wantErr    error
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
			wantValid: false,
			wantResult: SenderCSV{
				AccountId:      "",
				Address:        "",
				Country:        "AU",
				Channels:       []string{"sms"},
				MMSProviderKey: "",
				Comment:        "",
				Status:         "skipped",
				Error:          `Field "address" cannot be empty`,
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
			wantValid: false,
			wantResult: SenderCSV{
				AccountId:      "",
				Address:        "RHINO",
				Country:        "AU",
				Channels:       []string{"sms"},
				MMSProviderKey: "",
				Comment:        "",
				Status:         "skipped",
				Error:          `Field "address" must be unique`,
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
			valid, result, err := s.validateCSVSender(ctx, tt.args.csvSender)
			if err != nil {
				if tt.wantErr != nil {
					t.Errorf("senderImpl.validateCSVSender() error = %v, \n\twantErr %v", err, tt.wantErr)
				} else {
					t.Errorf("unexpected error: %+v", err)
				}
			}

			if valid != tt.wantValid {
				t.Errorf("senderImpl.validateCSVSender() valid = %v, \n\twant %v", valid, tt.wantValid)
			}
			if !reflect.DeepEqual(result, tt.wantResult) {
				t.Errorf("senderImpl.validateCSVSender() result = %+v, \n\twant %+v", result, tt.wantResult)
			}
		})
	}
}
