package tecloo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"text/template"

	"github.com/burstsms/mtmo-tp/backend/lib/mm7utils"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"gotest.tools/assert"
)

var (
	statusCodeRegex = regexp.MustCompile(`<StatusCode>(.*)<\/StatusCode>`)
	statusTextRegex = regexp.MustCompile(`<StatusText>(.*)<\/StatusText>`)
)

func TestSubmit(t *testing.T) {

	// Metadata content
	soaptmpl := template.Must(template.ParseFiles(`templates/tecloo_submit.soap.tmpl`))

	msgtext := `my test mms`

	image1, err := ioutil.ReadFile("image1.png")
	if err != nil {
		t.Error(err)
	}

	api := NewTeclooAPI(&TeclooAPIOptions{TemplatePath: `templates`})
	loghook := test.NewLocal(api.log.Logger)
	// loop over all the possible status codes and use to generate the recipient number
	for code, text := range statusCodes {

		// Request Content-Type with boundary parameter.
		submit := mm7utils.SubmitParams{
			TransactionID:    code,
			Subject:          "My MMS test message",
			VASPID:           "MYVASPID",
			Sender:           "61455123456",
			Recipient:        fmt.Sprintf("6142226%s", code),
			AllowAdaptations: true,
		}
		body, contentType, _ := mm7utils.GenerateMM7Submit(submit, soaptmpl, msgtext, [][]byte{image1})
		req, err := http.NewRequest("POST", "/v1/mm7", bytes.NewReader(body.Bytes()))
		req.Header.Set("Content-Type", contentType)
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

		// Parse the soap body
		soap := strings.Replace(rr.Body.String(), "\n", "", -1)
		soap = stripRegex.ReplaceAllString(soap, "><")

		transactionid := mm7utils.ExtractEntity(*transactionIDRegex, soap)
		assert.Equal(t, code, transactionid)

		statuscode := mm7utils.ExtractEntity(*statusCodeRegex, soap)
		assert.Equal(t, code, statuscode)

		statustext := mm7utils.ExtractEntity(*statusTextRegex, soap)
		assert.Equal(t, text.text, statustext)

		logentry := loghook.LastEntry()
		if statuscode == "1000" {
			drlogentry := logentry
			// if we were successful there will be a DR log entry also
			logentry = loghook.AllEntries()[(len(loghook.Entries) - 2)]
			log.Printf("%+v", loghook.Entries)
			assert.Equal(t, "MM7 DR", drlogentry.Message)
			assert.Equal(t, "1000", drlogentry.Data[`Status`])
			assert.Equal(t, drStatusCodes["1000"], drlogentry.Data[`StatusText`])
			assert.Equal(t, submit.Sender, drlogentry.Data[`Sender`])
			assert.Equal(t, submit.Recipient, drlogentry.Data[`Recipient`])
			assert.Equal(t, submit.TransactionID, drlogentry.Data[`TransactionID`])
		}
		// Make sure the logs expected
		assert.Equal(t, "MM7 Submit", logentry.Message)
		assert.Equal(t, logrus.InfoLevel, logentry.Level)
		assert.Equal(t, code, logentry.Data[`Status`])
		assert.Equal(t, submit.Subject, logentry.Data[`Subject`])
		assert.Equal(t, submit.Recipient, logentry.Data[`Recipient`])
		assert.Equal(t, submit.Sender, logentry.Data[`Sender`])
		assert.Equal(t, submit.TransactionID, logentry.Data[`TransactionID`])
		assert.Equal(t, msgtext, logentry.Data[`Message`])

		loghook.Reset()
	}
}
