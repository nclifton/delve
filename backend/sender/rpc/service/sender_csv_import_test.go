package service

import (
	"reflect"
	"testing"

	"github.com/vincent-petithory/dataurl"
)

func Test_unmarshalSenderCSVDataUrl(t *testing.T) {
	type args struct {
		csvDataUrl []byte
	}
	tests := []struct {
		name           string
		args           args
		wantCsvSenders []SenderCSV
		wantErr        bool
	}{
		{
			name: "CSV with two channels",
			args: args{
				csvDataUrl: []byte(
					dataurl.New(
						[]byte(`account_id,address,country,channels,mms_provider_key,comment
								,LION,AU,"[""mms"",""sms""]",fake,"roars"`),
						"text/csv").String()),
			},
			wantCsvSenders: []SenderCSV{{
				AccountId:      "",
				Address:        "LION",
				Country:        "AU",
				Channels:       []string{"mms", "sms"},
				MMSProviderKey: "fake",
				Comment:        "roars",
			}},
			wantErr: false,
		},
		{
			name: "CSV with one channel",
			args: args{
				csvDataUrl: []byte(
					dataurl.New(
						[]byte(`account_id,address,country,channels,mms_provider_key,comment
								,LION,AU,"[""mms""]",fake,"roars"`),
						"text/csv").String()),
			},
			wantCsvSenders: []SenderCSV{{
				AccountId:      "",
				Address:        "LION",
				Country:        "AU",
				Channels:       []string{"mms"},
				MMSProviderKey: "fake",
				Comment:        "roars",
			}},
			wantErr: false,
		}, {
			name: "CSV with no channels",
			args: args{
				csvDataUrl: []byte(
					dataurl.New(
						[]byte(`account_id,address,country,channels,mms_provider_key,comment
								,LION,AU,"[]",fake,"roars"`),
						"text/csv").String()),
			},
			wantCsvSenders: []SenderCSV{{
				AccountId:      "",
				Address:        "LION",
				Country:        "AU",
				Channels:       []string{},
				MMSProviderKey: "fake",
				Comment:        "roars",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCsvSenders, err := unmarshalSenderCSVDataUrl(tt.args.csvDataUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalSenderCSVDataUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCsvSenders, tt.wantCsvSenders) {
				t.Errorf("unmarshalSenderCSVDataUrl() = %v, want %v", gotCsvSenders, tt.wantCsvSenders)
			}
		})
	}
}
