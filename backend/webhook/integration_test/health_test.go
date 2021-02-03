// +build integration

package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HealthCheck(t *testing.T) {
	log.Println("test Webhook Health Check")

	s := setupForHealthCheckTest(t)
	defer s.teardown(t)

	client := http.DefaultClient

	type responseBody struct {
		GoRoutineThreshold string `json:"goroutine-threshold,omitempty"`
		Database           string `json:"database,omitempty"`
		Service            string `json:"service,omitempty"`
	}

	type want struct {
		httpResponseCode int
		responseBody     responseBody
	}

	tests := []struct {
		name         string
		path         string
		service      string
		stopPostgres bool
		pauseSeconds int
		want         want
		wantErr      error
	}{
		{
			name:    "RPC live",
			service: "rpc",
			path:    "/live?full=1",
			want: want{
				httpResponseCode: 200,
				responseBody: responseBody{
					GoRoutineThreshold: "OK",
				},
			},
		},
		{
			name:         "RPC ready database unavailable",
			service:      "rpc",
			path:         "/ready?full=1",
			stopPostgres: true,
			want: want{
				httpResponseCode: 200,
				responseBody: responseBody{
					GoRoutineThreshold: "OK",
					Database:           "^[^O][^K]",
					Service:            "OK",
				},
			},
		},
		{
			name:    "RPC ready",
			service: "rpc",
			path:    "/ready?full=1",
			want: want{
				httpResponseCode: 200,
				responseBody: responseBody{
					GoRoutineThreshold: "OK",
					Database:           "OK",
					Service:            "OK",
				},
			},
		},
		{
			name:         "worker live",
			service:      "worker",
			path:         "/ready?full=1",
			stopPostgres: false,
			want: want{
				httpResponseCode: 200,
				responseBody: responseBody{
					GoRoutineThreshold: "OK",
					Service:            "OK",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.stopPostgres {
				s.Postgres.Stop()
				// we need to teardown and restart all the fixtures now.
				// The postgres container has to be new after we stop it.
				// So to get the dependants on postgres to use a new postgres service they need to be new as well
				// Sorry.
				// Also note that for this need to be able to shut everything down without an os.Exit() happening
				defer func() {
					tfx.Teardown()
					setupFixtures()
					s.TestFixtures = *tfx
				}()
			}
			uri := s.RPCHealthCheckURI
			if tt.service == "worker" {
				uri = s.WorkerHealthCheckURIs[0]
			}
			url := fmt.Sprintf("%s%s", uri, tt.path)
			// t.Logf("GET: %s", url)
			resp, err := client.Get(url)
			if tt.wantErr != nil && err != nil {
				assert.Equal(t, tt.wantErr, err, "error")
			} else if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			// t.Logf("received: \n%s\n", string(bodyBytes))
			body := responseBody{}
			err = json.Unmarshal(bodyBytes, &body)
			assert.Regexp(t, tt.want.responseBody.GoRoutineThreshold, body.GoRoutineThreshold, "GoRoutineThreshold")
			assert.Regexp(t, tt.want.responseBody.Database, body.Database, "Database")
			assert.Regexp(t, tt.want.responseBody.Service, body.Service, "Service")
		})
	}
}

func setupForHealthCheckTest(t *testing.T) *testDeps {
	s := newSetup(t, tfx)

	return s
}
