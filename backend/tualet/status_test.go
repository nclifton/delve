package tualet

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStatus(t *testing.T) {

	api := NewTualetAPI(&TualetAPIOptions{TemplatePath: `templates`})
	req, err := http.NewRequest("GET", "/v1/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	api.router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"status":"ok"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'",
			rr.Body.String(), expected)
	}
}
