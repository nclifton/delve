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

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vincent-petithory/dataurl"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
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

	type want struct {
		createSendersParams *senderpb.CreateSendersParams
		createSendersReply  *senderpb.CreateSendersReply
		createSendersError  error
		bodyString          string
		statusCode          int
	}

	tests := []struct {
		name    string
		csv     []string
		want    want
	}{
		{
			name: "happy import",
			csv: []string{
				"account_id,address,country,channels,mms_provider_key,comment",
				`,GIRAFFE,AU,"[""sms"",""mms""]",,`,
				`,NOKEY,AU,"[""sms""]",,`,
			},
			want: want{
				createSendersParams: &senderpb.CreateSendersParams{
					Senders: []*senderpb.NewSender{
						{
							AccountId:      "",
							Address:        "GIRAFFE",
							MMSProviderKey: "",
							Channels:       []string{"sms", "mms"},
							Country:        "AU",
							Comment:        "",
						},
						{
							AccountId:      "",
							Address:        "NOKEY",
							MMSProviderKey: "",
							Channels:       []string{"sms"},
							Country:        "AU",
							Comment:        "",
						},
					},
				},
				createSendersReply: &senderpb.CreateSendersReply{
					Senders: []*senderpb.Sender{
						{
							Id:        uuid.New(),
							Address:   "GIRAFFE",
							Channels:  []string{"sms", "mms"},
							Country:   "AU",
							Comment:   "",
							CreatedAt: timestamppb.Now(),
							UpdatedAt: timestamppb.Now(),
						},
						{
							Id:        uuid.New(),
							Address:   "NOKEY",
							Channels:  []string{"sms"},
							Country:   "AU",
							Comment:   "",
							CreatedAt: timestamppb.Now(),
							UpdatedAt: timestamppb.Now(),
						},
					},
				},
				bodyString: `{"status":"ok"}`,
				statusCode: http.StatusOK,
			},
		}, {
			name: "unhappy - address not provided",
			csv: []string{
				"account_id,address,country,channels,mms_provider_key,comment",
				`,,AU,"[""sms"",""mms""]",,`,
			},
			want: want{
				createSendersParams: &senderpb.CreateSendersParams{
					Senders: []*senderpb.NewSender{{
						AccountId:      "",
						Address:        "",
						MMSProviderKey: "",
						Channels:       []string{"sms", "mms"},
						Country:        "AU",
						Comment:        "",
					}},
				},
				createSendersReply: &senderpb.CreateSendersReply{},
				createSendersError: status.Error(codes.Unknown, `ERROR: null value in column "address" violates not-null constraint (SQLSTATE 23502)`),
				bodyString:         `{"error":"Could not upload senders CSV: ERROR: null value in column \"address\" violates not-null constraint (SQLSTATE 23502)"}`,
				statusCode:         http.StatusInternalServerError,
			},
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

			// prepare and inject the mock sender RPC service client
			mock := new(senderpb.MockServiceClient)
			mock.On("CreateSenders", req.Context(), tt.want.createSendersParams).Return(tt.want.createSendersReply, tt.want.createSendersError)
			api := NewAdminAPI(&AdminAPIOptions{
				SenderClient: mock,
			})

			rr := httptest.NewRecorder()
			api.Handler().ServeHTTP(rr, req)

			// Check the status code is what we expect.
			assert.Equal(t, tt.want.statusCode, rr.Code, "handler returned wrong status code")

			// Check the response body is what we expect.
			assert.JSONEq(t, tt.want.bodyString, rr.Body.String(), "handler returned unexpected body")

		})
	}
}
