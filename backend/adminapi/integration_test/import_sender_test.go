// +build integration

package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/burstsms/mtmo-tp/backend/adminapi"
)

func Test_ImportSenderPOST(t *testing.T) {

	api := adminapi.NewAdminAPI(&adminapi.AdminAPIOptions{})
	req, err := http.NewRequest("GET", "/v1/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	api.Handler().ServeHTTP(rr, req)

	


}