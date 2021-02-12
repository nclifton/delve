package adminapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vincent-petithory/dataurl"
)

func TestArray_UnmarshalCSV(t *testing.T) {
	type args struct {
		csv string
	}
	tests := []struct {
		name    string
		a       *Array
		want    *Array
		args    args
		wantErr bool
	}{
		{
			name: "two channels",
			a:    &Array{},
			args: args{
				csv: `["mms","sms"]`,
			},
			want:    &Array{"mms", "sms"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.UnmarshalCSV(tt.args.csv); (err != nil) != tt.wantErr {
				t.Errorf("Array.UnmarshalCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, tt.a)
		})
	}
}

func TestGetSendersFromRequest(t *testing.T) {

	tests := []struct {
		name   string
		csv    []string
		want   []SenderCSV
		failed bool
	}{
		{
			name: "happy two senders one with two channels",
			csv: []string{
				"account_id,address,country,channels,mms_provider_key,comment",
				`,GIRAFFE,AU,"[""sms"",""mms""]",,`,
				`,NOKEY,AU,"[""sms""]",,`,
			},
			want: []SenderCSV{
				{
					AccountId:      "",
					Address:        "GIRAFFE",
					Country:        "AU",
					Channels:       []string{"sms", "mms"},
					MMSProviderKey: "",
					Comment:        "",
				}, {
					AccountId:      "",
					Address:        "NOKEY",
					Country:        "AU",
					Channels:       []string{"sms"},
					MMSProviderKey: "",
					Comment:        "",
				}},
			failed: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j, err := json.Marshal(ImportJSON{
				Data: dataurl.New([]byte(strings.Join(tt.csv, "\n")), "text/csv").String(),
			})
			if err != nil {
				t.Fatal(err)
			}
			route := &Route{
				r: &http.Request{
					Body:     ioutil.NopCloser(bytes.NewReader(j)),
					Response: &http.Response{},
				},
			}
			got, failed := GetSendersFromRequest(route)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSendersFromRequest() got = %+v, want %+v", got, tt.want)
			}
			if failed != tt.failed {
				t.Errorf("GetSendersFromRequest() got1 = %+v, want %+v", failed, tt.failed)
			}
		})
	}
}

func Test_ImportSenderPOST(t *testing.T) {

	api := NewAdminAPI(&AdminAPIOptions{})

	type wantErr error
	type want struct {
		bodyString string
		statusCode int
	}

	tests := []struct {
		name    string
		csv     []string
		want    want
		wantErr wantErr
	}{
		{
			name: "happy import",
			csv: []string{
				"account_id,address,country,channels,mms_provider_key,comment",
				`,GIRAFFE,AU,"[""sms"",""mms""]",,`,
				`,NOKEY,AU,"[""sms""]",,`,
			},
			want: want{
				bodyString: `{"status":"ok"}`,
				statusCode: http.StatusOK,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			j, err := json.Marshal(ImportJSON{
				Data: dataurl.New([]byte(strings.Join(tt.csv, "\n")), "text/csv").String(),
			})
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest("POST", "/v1/import/sender", bytes.NewBuffer(j))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			api.Handler().ServeHTTP(rr, req)

			// Check the status code is what we expect.
			assert.Equal(t, tt.want.statusCode, rr.Code, "handler returned wrong status code")

			// Check the response body is what we expect.
			assert.JSONEq(t, tt.want.bodyString, rr.Body.String(), "handler returned unexpected body")

			// check the sender rpc client mock that it was used as exected

		})
	}
}
