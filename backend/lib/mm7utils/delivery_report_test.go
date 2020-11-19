package mm7utils

import (
	"fmt"
	"testing"
	"text/template"

	"gotest.tools/assert"
)

func TestGenerateMM7DeliveryReport(t *testing.T) {
	params := DeliveryReportParams{
		StatusCode:    "1000",
		StatusText:    "statustext",
		TransactionID: "2000",
		MessageID:     "3000",
		Recipient:     "61422265707",
		Sender:        "61422265777",
	}

	soaptmpl, err := template.New("test").Parse(`[{{.StatusCode}}],[{{.StatusText}}],[{{.TransactionID}}],[{{.MessageID}}],[{{.Sender}}],[{{.Recipient}}]`)
	if err != nil {
		t.Error(err)
	}

	soapdata, contentType, err := GenerateMM7DeliveryReport(params, soaptmpl)
	if err != nil {
		t.Error(err)
	}

	assert.Assert(t, testContains(contentType, "application/xml"))
	assert.Assert(t, testContains(soapdata.String(), fmt.Sprintf(`[%s]`, params.StatusCode)))
	assert.Assert(t, testContains(soapdata.String(), fmt.Sprintf(`[%s]`, params.StatusText)))
	assert.Assert(t, testContains(soapdata.String(), fmt.Sprintf(`[%s]`, params.TransactionID)))
	assert.Assert(t, testContains(soapdata.String(), fmt.Sprintf(`[%s]`, params.MessageID)))
	assert.Assert(t, testContains(soapdata.String(), fmt.Sprintf(`[%s]`, params.Sender)))
	assert.Assert(t, testContains(soapdata.String(), fmt.Sprintf(`[%s]`, params.Recipient)))

}

func TestGenerateMM7DeliverReportResponse(t *testing.T) {
	params := DeliveryReportResponseParams{
		TransactionID: "10001",
		StatusCode:    "1000",
		StatusText:    "statustext",
	}

	soaptmpl, err := template.New("test").Parse(`[{{.TransactionID}}],[{{.StatusCode}}],[{{.StatusText}}]`)
	if err != nil {
		t.Error(err)
	}

	soapdata, contentType, err := GenerateMM7DeliveryReportResponse(params, soaptmpl)
	if err != nil {
		t.Error(err)
	}

	assert.Assert(t, testContains(contentType, "text/xml"))
	assert.Assert(t, testContains(soapdata.String(), fmt.Sprintf(`[%s]`, params.TransactionID)))
	assert.Assert(t, testContains(soapdata.String(), fmt.Sprintf(`[%s]`, params.StatusCode)))
	assert.Assert(t, testContains(soapdata.String(), fmt.Sprintf(`[%s]`, params.StatusText)))
}
