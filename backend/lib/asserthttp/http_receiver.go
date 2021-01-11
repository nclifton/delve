package asserthttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type HttpReceiver struct {
	httpServer        *httptest.Server
	httpRequests      []*http.Request
	httpRequestBodies []string
	t                 *testing.T
}

func NewHttpReceiver(t *testing.T) *HttpReceiver {
	receiver := &HttpReceiver{
		t:                 t,
		httpRequests:      make([]*http.Request, 0),
		httpRequestBodies: make([]string, 0),
	}
	receiver.newServer()
	log.Printf("started receiver at %s", receiver.httpServer.URL)
	return receiver
}
func (receiver *HttpReceiver) newServer() {
	receiver.httpServer = httptest.NewServer(http.HandlerFunc(receiver.handler))
}

func (receiver *HttpReceiver) handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		receiver.t.Fatalf("failed to read request body, %+v", err)
	}

	// should I use a channel to record the received requests?
	receiver.httpRequests = append(receiver.httpRequests, r)
	receiver.httpRequestBodies = append(receiver.httpRequestBodies, string(body))

	fmt.Fprintln(w, "thank you")
	log.Printf("received http request at %s", receiver.httpServer.URL)
}

func (s *HttpReceiver) GetURL() string {
	return s.httpServer.URL
}

func (s *HttpReceiver) Teardown() {
	log.Printf("closing http receiver at %s", s.httpServer.URL)
	s.httpServer.Close()
}

type ExpectedRequests struct {
	NumberOfRequests int
	WaitMilliseconds int
	Methods          []string
	ContentTypes     []string
	Bodies           []string
}

func (s *HttpReceiver) WaitForRequests(want ExpectedRequests) {
	var cnt = 0
	log.Printf("waiting for http request at %s", s.httpServer.URL)
	for len(s.httpRequests) < want.NumberOfRequests || want.NumberOfRequests == 0 {
		if cnt > want.WaitMilliseconds {
			plural := ""
			if want.NumberOfRequests > 1 {
				plural = "s"
			}
			if want.NumberOfRequests > len(s.httpRequests) {
				assert.Fail(s.t, fmt.Sprintf("timed out waiting for %d request%s at %s", want.NumberOfRequests, plural, s.httpServer.URL))
				return
			} else {
				log.Printf("http request wait timed out")
				return
			}
		}
		// TODO use a waitgroup or a channel instead of sleep polling
		time.Sleep(time.Millisecond)
		cnt++
	}
	log.Printf("received http request at %s after %d milliseconds", s.httpServer.URL, cnt)
}

func (s *HttpReceiver) ResetHttpRequests() {
	s.httpRequests = nil
	s.httpRequestBodies = nil
}

func JSONString(t *testing.T, v interface{}) string {
	expectedBody, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("error: %+v", err)
	}
	return string(expectedBody)
}

func (s *HttpReceiver) AssertRequests(want ExpectedRequests) {
	if assert.Equal(s.t, want.NumberOfRequests, len(s.httpRequests), "number of requests received") {
		for i, method := range want.Methods {
			if i < len(s.httpRequests) {
				req := s.httpRequests[i]
				assert.Equal(s.t, req.Method, method, fmt.Sprintf("request %d method", i+1))
			}
		}
		for i, contentType := range want.ContentTypes {
			if i < len(s.httpRequests) {
				req := s.httpRequests[i]
				assert.Equal(s.t, req.Header.Get("Content-Type"), contentType, fmt.Sprintf("request %d Content-Type", i+1))
			}
		}
		for i, body := range want.Bodies {
			if i < len(s.httpRequestBodies) {
				assert.JSONEq(s.t, body, s.httpRequestBodies[i], fmt.Sprintf("request %d body json", i+1))
			}
		}
	}
}
