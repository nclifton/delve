package mm7utils

import (
	"fmt"
	"testing"
	"text/template"

	"gotest.tools/assert"
)

func TestGenerateMM7SubmitResponse(t *testing.T) {
	params := SubmitResponseParams{
		StatusCode:    "1000",
		StatusText:    "statustext",
		TransactionID: "2000",
		MessageID:     "3000",
	}

	soaptmpl, err := template.New("test").Parse(`[{{.StatusCode}}],[{{.StatusText}}],[{{.TransactionID}}],[{{.MessageID}}]`)
	if err != nil {
		t.Error(err)
	}

	soapdata, contentType, err := GenerateMM7SubmitResponse(params, soaptmpl)
	if err != nil {
		t.Error(err)
	}

	assert.Assert(t, testContains(contentType, "application/xml"))
	assert.Assert(t, testContains(soapdata.String(), fmt.Sprintf(`[%s]`, params.StatusCode)))
	assert.Assert(t, testContains(soapdata.String(), fmt.Sprintf(`[%s]`, params.StatusText)))
	assert.Assert(t, testContains(soapdata.String(), fmt.Sprintf(`[%s]`, params.TransactionID)))
	assert.Assert(t, testContains(soapdata.String(), fmt.Sprintf(`[%s]`, params.MessageID)))

}
