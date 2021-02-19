package adminapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vincent-petithory/dataurl"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
)

func Test_ImportSenderPOST(t *testing.T) {

	type mockStuff struct {
		createSendersReply *senderpb.CreateSendersFromCSVDataURLReply
		createSendersError error
	}

	type want struct {
		mockStuff  mockStuff
		bodyString string
		statusCode int
		jsonError  *JSONErrors
	}

	tests := []struct {
		name string
		csv  string
		want want
	}{
		{
			name: "happy import",
			csv: `account_id,address,country,channels,mms_provider_key,comment
				,GIRAFFE,AU,"[""sms"",""mms""]",,
				,NOKEY,AU,"[""sms""]",,`,
			want: want{
				mockStuff: mockStuff{
					createSendersReply: &senderpb.CreateSendersFromCSVDataURLReply{
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
				},
				bodyString: `{"status":"ok"}`,
				statusCode: http.StatusOK,
			},
		}, {
			name: "unvalidated CSV - some unknown internal error",
			csv: `account_id,address,country,channels,mms_provider_key,comment
				,,AU,"[""sms"",""mms""]",,`,
			want: want{
				mockStuff: mockStuff{
					createSendersReply: &senderpb.CreateSendersFromCSVDataURLReply{},
					createSendersError: status.Error(codes.Unknown, `something bad happened`),
				},
				bodyString: `{"error":"Could not upload senders CSV: something bad happened"}`,
				statusCode: http.StatusInternalServerError,
			},
		}, {
			name: "empty data",
			csv:  "",
			want: want{
				statusCode: http.StatusUnprocessableEntity,
				jsonError: &JSONErrors{
					Error:     "Validation Error",
					ErrorData: map[string]string{"data": "required"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var duBytes []byte
			var err error
			if tt.csv != "" {
				duBytes, err = dataurl.New([]byte(tt.csv), "text/csv").MarshalText()
				if err != nil {
					t.Fatal(err)
				}
			}
			j, err := json.Marshal(ImportSenderPOSTRequest{
				Data: string(duBytes),
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
			params := &senderpb.CreateSendersFromCSVDataURLParams{CSV: post.Data}
			mock := new(senderpb.MockServiceClient)
			mock.On("CreateSendersFromCSVDataURL", req.Context(), params).Return(tt.want.mockStuff.createSendersReply, tt.want.mockStuff.createSendersError)
			api := NewAdminAPI(&AdminAPIOptions{
				SenderClient: mock,
			})

			api.Handler().ServeHTTP(rr, req)

			// Check the status code is what we expect.
			assert.Equal(t, tt.want.statusCode, rr.Code, "handler returned wrong status code")

			// Check the response body is what we expect.
			if tt.want.jsonError == nil {
				assert.JSONEq(t, tt.want.bodyString, rr.Body.String(), "handler returned unexpected body")
			} else {
				bytes, err := json.Marshal(tt.want.jsonError)
				if err != nil {
					t.Fatalf("wanted JSON Errors failed to marshal %+v", err)
				}
				assert.JSONEq(t, string(bytes), rr.Body.String(), "handler returned unexpected body")
			}

		})
	}
}
