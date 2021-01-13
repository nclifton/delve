package tecloo

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"gotest.tools/assert"
	"gotest.tools/assert/cmp"
)

func testContains(haystack string, needle string) cmp.Comparison {
	return func() cmp.Result {
		if strings.Contains(haystack, needle) {
			return cmp.ResultSuccess
		}
		return cmp.ResultFailure(fmt.Sprintf(`%s did not contain %s`, haystack, needle))
	}
}

func TestDRSend(t *testing.T) {

	// loop over all the possible status codes and use to generate the recipient number
	for code, text := range drStatusCodes {

		submit := DRParams{
			TransactionID: "1000",
			MessageID:     "1001",
			Sender:        "61455123456",
			Recipient:     fmt.Sprintf("6142226%s", code),
			Status:        code,
		}

		testDRhandler := func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Error(err)
			}
			req := string(body)
			assert.Assert(t, testContains(r.Header.Get(`Content-Type`), "application/xml"))
			assert.Assert(t, testContains(req, fmt.Sprintf(`<MMStatus>%s</MMStatus>`, submit.Status)))
			assert.Assert(t, testContains(req, fmt.Sprintf(`<StatusText>%s</StatusText>`, text)))
			assert.Assert(t, testContains(req, fmt.Sprintf(`%s</TransactionID>`, submit.TransactionID)))
			assert.Assert(t, testContains(req, fmt.Sprintf(`<MessageID>%s</MessageID>`, submit.MessageID)))
			assert.Assert(t, testContains(req, fmt.Sprintf(`<Number>%s</Number>`, submit.Sender)))
			assert.Assert(t, testContains(req, fmt.Sprintf(`<Number>%s</Number>`, submit.Recipient)))

			if r.Method != "POST" {
				t.Errorf("want %s, got %s", "GET", r.Method)
			}
		}
		s := httptest.NewServer(http.HandlerFunc(testDRhandler))

		api := NewTeclooAPI(&TeclooAPIOptions{TemplatePath: `templates`, DREndpoint: s.URL, Client: s.Client()})
		loghook := test.NewLocal(api.log.Logger)

		api.sendDRRequest(context.Background(), &submit)

		logentry := loghook.AllEntries()[len(loghook.Entries)-1]
		assert.Equal(t, "MM7 DR", logentry.Message)
		assert.Equal(t, code, logentry.Data[`Status`])
		assert.Equal(t, text, logentry.Data[`StatusText`])
		assert.Equal(t, submit.Sender, logentry.Data[`Sender`])
		assert.Equal(t, submit.Recipient, logentry.Data[`Recipient`])
		assert.Equal(t, submit.TransactionID, logentry.Data[`TransactionID`])
		assert.Equal(t, submit.MessageID, logentry.Data[`MessageID`])

		loghook.Reset()
	}
}
