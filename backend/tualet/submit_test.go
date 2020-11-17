package tualet

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestSubmit(t *testing.T) {

	api := NewTualetAPI(&TualetAPIOptions{TemplatePath: `templates`})
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := url.Values{}
	q.Add("username", "froggy")
	q.Add("password", "snoggy")
	q.Add("command", "submit")
	q.Add("message", "testing")
	q.Add("dnis", "61455123456")

	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()

	api.router.ServeHTTP(rr, req)

	log.Printf("req: %s", req.URL.String())

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"message_id":"xxx"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'",
			rr.Body.String(), expected)
	}
}
