package tualet

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"gotest.tools/assert"
)

func TestSubmit(t *testing.T) {

	testValues := []struct {
		values               submitParams
		expectedResponseCode int
		expectedResponseBody string
		expectedLogMessage   string
	}{
		{
			values: submitParams{
				username: "froggy",
			},
			expectedResponseCode: http.StatusUnauthorized,
			expectedResponseBody: `not authorized (check login and password)`,
			expectedLogMessage:   `not authorized (check login and password)`,
		},
		{
			values: submitParams{
				username:        "froggy",
				password:        "snoggy",
				command:         "submit",
				message:         "testing",
				dnis:            "61455123456",
				ani:             "61455678900",
				longMessageMode: "split",
			},
			expectedResponseCode: http.StatusOK,
			expectedResponseBody: `{"message_id":"`,
			expectedLogMessage:   `submission`,
		},
		{
			values: submitParams{
				username:        "froggy",
				password:        "snoggy",
				command:         "submit",
				message:         "testing",
				dnis:            "61455121400",
				ani:             "61455678900",
				longMessageMode: "split",
			},
			expectedResponseCode: http.StatusBadRequest,
			expectedResponseBody: `NO ROUTES`,
			expectedLogMessage:   `NO ROUTES`,
		},
	}

	api := NewTualetAPI(&TualetAPIOptions{})
	loghook := test.NewLocal(api.log.Logger)
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, testValue := range testValues {
		q := url.Values{}
		if testValue.values.username != "" {
			q.Add("username", testValue.values.username)
		}
		if testValue.values.password != "" {
			q.Add("password", testValue.values.password)
		}
		if testValue.values.command != "" {
			q.Add("command", testValue.values.command)
		}
		if testValue.values.message != "" {
			q.Add("message", testValue.values.message)
		}
		if testValue.values.dnis != "" {
			q.Add("dnis", testValue.values.dnis)
		}
		if testValue.values.ani != "" {
			q.Add("ani", testValue.values.ani)
		}
		if testValue.values.longMessageMode != "" {
			q.Add("longMessageMode", testValue.values.longMessageMode)
		}

		req.URL.RawQuery = q.Encode()

		rr := httptest.NewRecorder()

		api.router.ServeHTTP(rr, req)

		log.Printf("Req URL: %s", req.URL.String())

		// Check the status code is what we expect.
		if status := rr.Code; status != testValue.expectedResponseCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check the response body is what we expect.
		if !strings.Contains(strings.TrimSpace(rr.Body.String()), testValue.expectedResponseBody) {
			t.Errorf("handler returned unexpected body: got '%v' want '%v'",
				rr.Body.String(), testValue.expectedResponseBody)
		}

		logentry := loghook.LastEntry()
		assert.Equal(t, testValue.expectedLogMessage, logentry.Message)
		assert.Equal(t, logrus.InfoLevel, logentry.Level)
		assert.Equal(t, testValue.values.message, logentry.Data[`message`])
		assert.Equal(t, testValue.values.command, logentry.Data[`command`])
		assert.Equal(t, testValue.values.dnis, logentry.Data[`dnis`])
		assert.Equal(t, testValue.values.ani, logentry.Data[`ani`])
		assert.Equal(t, testValue.values.longMessageMode, logentry.Data[`longMessageMode`])
		loghook.Reset()
	}
}
