package tecloo_receiver

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/burstsms/mtmo-tp/backend/lib/mm7utils"
	"gotest.tools/assert"
)

func TestInboundPOST(t *testing.T) {
	var (
		deliverRspRegex        = regexp.MustCompile(`<DeliverRsp.*>`)
		deliveryReportRspRegex = regexp.MustCompile(`<DeliveryReportRsp.*>`)
		drsoaptmpl             = template.Must(template.ParseFiles(`templates/mm7_deliver.soap.tmpl`))
		dlrsoaptmpl            = template.Must(template.ParseFiles(`templates/mm7_delivery_report.soap.tmpl`))
	)

	api := NewTeclooReceiverAPI(&TeclooReceiverAPIOptions{
		TemplatePath: `templates`,
	})

	t.Run("valid delivery request results in success delivery response", func(t *testing.T) {
		deliverRequest := mm7utils.DeliverRequestParams{
			TransactionID: "1000",
			Sender:        "61455123456",
			Recipient:     "111122",
			Date:          time.Now().UTC(),
		}
		image1, err := ioutil.ReadFile("image1.png")
		if err != nil {
			t.Error(err)
		}
		body, contentType, _ := mm7utils.GenerateMM7DeliverRequest(deliverRequest, drsoaptmpl, "inbound success", [][]byte{image1})
		req, err := http.NewRequest("POST", "/v1/mms/inbound", bytes.NewReader(body.Bytes()))
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

		assert.Equal(t, true, deliverRspRegex.MatchString(soap))

		statusCode := extractEntity(*statusCodeRegex, soap)
		assert.Equal(t, "1000", statusCode)

		statusText := extractEntity(*statusTextRegex, soap)
		assert.Equal(t, mm7utils.StatusCodes["1000"].Text, statusText)
	})

	t.Run("invalid delivery request results in error delivery response", func(t *testing.T) {
		deliverRequest := mm7utils.DeliverRequestParams{
			TransactionID: "1001",
			Sender:        "61455122001",
			Recipient:     "61455122007",
			Date:          time.Now().UTC(),
		}
		image1, err := ioutil.ReadFile("image1.png")
		if err != nil {
			t.Error(err)
		}
		body, contentType, _ := mm7utils.GenerateMM7DeliverRequest(deliverRequest, drsoaptmpl, "inbound fail", [][]byte{image1})
		req, err := http.NewRequest("POST", "/v1/mms/inbound", bytes.NewReader(body.Bytes()))
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

		assert.Equal(t, true, deliverRspRegex.MatchString(soap))

		statusCode := extractEntity(*statusCodeRegex, soap)
		assert.Equal(t, "2007", statusCode)

		statusText := extractEntity(*statusTextRegex, soap)
		assert.Equal(t, mm7utils.StatusCodes["2007"].Text, statusText)
	})

	t.Run("valid delivery report request results in success delivery report response", func(t *testing.T) {
		deliveryReportRequest := mm7utils.DeliveryReportParams{
			TransactionID: "100001",
			StatusCode:    "1000",
			StatusText:    "success",
			Date:          time.Now().UTC().String(),
			MessageID:     "100001",
			Recipient:     "61455122001",
			Sender:        "61455122002",
		}

		body, contentType, _ := mm7utils.GenerateMM7DeliveryReport(deliveryReportRequest, dlrsoaptmpl)
		req, err := http.NewRequest("POST", "/v1/mms/inbound", bytes.NewReader(body.Bytes()))
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

		assert.Equal(t, true, deliveryReportRspRegex.MatchString(soap))

		statusCode := extractEntity(*statusCodeRegex, soap)
		assert.Equal(t, "1000", statusCode)

		statusText := extractEntity(*statusTextRegex, soap)
		assert.Equal(t, mm7utils.StatusCodes["1000"].Text, statusText)
	})

	t.Run("invalid delivery report request results in error delivery report response", func(t *testing.T) {
		deliveryReportRequest := mm7utils.DeliveryReportParams{
			TransactionID: "200001",
			StatusCode:    "2007",
			StatusText:    "Unable to parse request",
			Date:          time.Now().UTC().String(),
			MessageID:     "100001",
			Recipient:     "61455122007",
			Sender:        "61455122002",
		}

		body, contentType, _ := mm7utils.GenerateMM7DeliveryReport(deliveryReportRequest, dlrsoaptmpl)
		req, err := http.NewRequest("POST", "/v1/mms/inbound", bytes.NewReader(body.Bytes()))
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

		assert.Equal(t, true, deliveryReportRspRegex.MatchString(soap))

		statusCode := extractEntity(*statusCodeRegex, soap)
		assert.Equal(t, "2007", statusCode)

		statusText := extractEntity(*statusTextRegex, soap)
		assert.Equal(t, mm7utils.StatusCodes["2007"].Text, statusText)
	})
}
