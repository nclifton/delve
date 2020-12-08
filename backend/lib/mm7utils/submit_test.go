package mm7utils

import (
	"fmt"
	"io/ioutil"
	"testing"
	"text/template"

	"gotest.tools/assert"
)

func TestGenerateMM7Submit(t *testing.T) {
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

	submit := SubmitParams{
		TransactionID:    "1000",
		Subject:          "My MMS test message",
		VASPID:           "MYVASPID",
		Sender:           "61455123456",
		Recipient:        fmt.Sprintf("6142226%s", "1000"),
		AllowAdaptations: true,
	}

	body, contentType, _ := GenerateMM7Submit(submit, soaptmpl, msgtext, [][]byte{image1, image2})
	assert.Assert(t, testContains(contentType, "multipart/related"))

	parts, err := ProcessMultiPart(contentType, body)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, len(parts), 5)
	assert.Equal(t, parts[0].ContentID, "soap-start")
	assert.Equal(t, string(parts[0].Body), "1000")
	assert.Equal(t, parts[0].ContentType, "text/xml")
	assert.Equal(t, parts[1].ContentID, "msg-txt")
	assert.Equal(t, parts[1].ContentType, "text/plain")
	assert.Equal(t, string(parts[1].Body), msgtext)
	assert.Equal(t, parts[2].ContentID, "image-0")
	assert.Equal(t, parts[2].ContentType, "image/png")
	assert.Equal(t, parts[3].ContentID, "image-1")
	assert.Equal(t, parts[3].ContentType, "image/png")
	assert.Equal(t, parts[4].ContentID, "mms.smil")
	assert.Equal(t, parts[4].ContentType, "application/smil")

}
