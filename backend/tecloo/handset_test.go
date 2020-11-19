package tecloo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"gotest.tools/assert"
)

func TestHandset(t *testing.T) {

	mux := http.NewServeMux()

	image1, err := ioutil.ReadFile("image1.png")
	if err != nil {
		t.Error(err)
	}

	testImghandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(`Content-Type`, "image/png")
		w.Header().Set(`Content-Length`, strconv.Itoa(len(image1)))
		if _, err := w.Write(image1); err != nil {
			t.Error(err)
		}
	}

	mux.HandleFunc("/image1.png", testImghandler)

	testDRhandler := func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		assert.Assert(t, testContains(r.Header.Get(`Content-Type`), "multipart/related"))
		parts, err := ProcessMultiPart(r.Header.Get(`Content-Type`), bytes.NewReader(body))
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, len(parts), 4)
		assert.Equal(t, parts[0].ContentID, "<soap-start>")
		assert.Equal(t, parts[0].ContentType, "text/xml")
		assert.Equal(t, parts[1].ContentID, "<msg-txt>")
		assert.Equal(t, parts[1].ContentType, "text/plain")
		assert.Equal(t, parts[2].ContentID, "image-0.png")
		assert.Equal(t, parts[2].ContentType, "image/png")
		assert.Equal(t, parts[3].ContentID, "<mms.smil>")
		assert.Equal(t, parts[3].ContentType, "application/smil")

		if r.Method != "POST" {
			t.Errorf("want %s, got %s", "POST", r.Method)
		}
	}

	mux.HandleFunc("/dr", testDRhandler)

	s := httptest.NewServer(mux)

	api := NewTeclooAPI(&TeclooAPIOptions{TemplatePath: `templates`, Client: s.Client(), DREndpoint: s.URL + "/dr"})
	loghook := test.NewLocal(api.log.Logger)

	submit := HandsetParams{
		ID:          "vwefwefgwefwefwe",
		Subject:     "Handset Subject",
		Message:     "MMS from tecloo handset",
		Recipient:   "61455671000",
		Sender:      "61455679998",
		ProviderKey: "MYVASPID",
		ContentURLs: []string{s.URL + "/image1.png"},
	}

	jsonValue, err := json.Marshal(submit)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/v1/handset/mms", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
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

	loghook.Reset()
}
