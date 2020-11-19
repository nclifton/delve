package mm7utils

import (
	"fmt"
	"io/ioutil"
	"testing"
	"text/template"
	"time"

	"gotest.tools/assert"
)

func TestGenerateMM7DeliverRequest(t *testing.T) {
	msgtext := `my test mms`
	soaptmpl, err := template.New("test").Parse(`{{.TransactionID}}`)
	if err != nil {
		t.Error(err)
	}

	image1, err := ioutil.ReadFile("image1.png")
	if err != nil {
		t.Error(err)
	}
	image2, err := ioutil.ReadFile("image1.png")
	if err != nil {
		t.Error(err)
	}

	params := DeliverRequestParams{
		TransactionID: "transactionid",
		Sender:        "sender",
		Recipient:     "recipient",
		Date:          time.Now().UTC(),
	}

	soapdata, contentType, _ := GenerateMM7DeliverRequest(params, soaptmpl, msgtext, [][]byte{image1, image2})
	assert.Assert(t, testContains(contentType, "multipart/related"))

	parts, err := ProcessMultiPart(contentType, soapdata)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, len(parts), 5)
	assert.Equal(t, parts[0].ContentID, "<soap-start>")
	assert.Equal(t, string(parts[0].Body), "transactionid")
	assert.Equal(t, parts[0].ContentType, "text/xml")
	assert.Equal(t, parts[1].ContentID, "<msg-txt>")
	assert.Equal(t, parts[1].ContentType, "text/plain")
	assert.Equal(t, string(parts[1].Body), msgtext)
	assert.Equal(t, parts[2].ContentID, "image-0.png")
	assert.Equal(t, parts[2].ContentType, "image/png")
	assert.Equal(t, parts[3].ContentID, "image-1.png")
	assert.Equal(t, parts[3].ContentType, "image/png")
	assert.Equal(t, parts[4].ContentID, "<mms.smil>")
	assert.Equal(t, parts[4].ContentType, "application/smil")
}

func TestGenerateMM7DeliverResponse(t *testing.T) {
	params := DeliverResponseParams{
		TransactionID: "transactionid",
		StatusCode:    "1000",
		StatusText:    "statustext",
	}

	soaptmpl, err := template.New("test").Parse(`[{{.TransactionID}}],[{{.StatusCode}}],[{{.StatusText}}]`)
	if err != nil {
		t.Error(err)
	}

	soapdata, contentType, err := GenerateMM7DeliverResponse(params, soaptmpl)
	if err != nil {
		t.Error(err)
	}

	assert.Assert(t, testContains(contentType, "application/xml"))
	assert.Assert(t, testContains(soapdata.String(), fmt.Sprintf(`[%s]`, params.StatusCode)))
	assert.Assert(t, testContains(soapdata.String(), fmt.Sprintf(`[%s]`, params.StatusText)))
}
