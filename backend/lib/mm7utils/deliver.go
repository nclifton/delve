package mm7utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"text/template"
	"time"
)

type DeliverRequestParams struct {
	TransactionID    string
	Sender           string
	Recipient        string
	Subject          string
	ContentID        string
	Date             time.Time
	Priority         string
	AllowAdaptations bool
}

func GenerateMM7DeliverRequest(params DeliverRequestParams, soaptmpl *template.Template, message string, images [][]byte) (*bytes.Buffer, string, error) {
	// New multipart writer.
	body := &bytes.Buffer{}
	contentBody := &bytes.Buffer{}
	soapdata := &bytes.Buffer{}
	main := multipart.NewWriter(body)
	content := multipart.NewWriter(contentBody)

	// Metadata part.
	metadataHeader := textproto.MIMEHeader{}
	metadataHeader.Set("Content-Type", "text/xml")
	metadataHeader.Set("Content-ID", "<soap-start>")
	part, err := main.CreatePart(metadataHeader)
	if err != nil {
		return nil, "", err
	}

	params.ContentID = "mmscontent"
	err = soaptmpl.Execute(soapdata, params)
	if err != nil {
		return nil, "", err
	}
	_, err = part.Write(soapdata.Bytes())
	if err != nil {
		return nil, "", err
	}

	// Text Part
	txtHeader := textproto.MIMEHeader{}
	txtHeader.Set("Content-Type", "text/plain")
	txtHeader.Set("Content-ID", "<msg-txt>")
	txtPart, err := content.CreatePart(txtHeader)
	if err != nil {
		return nil, "", err
	}
	_, err = txtPart.Write([]byte(message))
	if err != nil {
		return nil, "", err
	}
	attachments := []SMILMedia{{ContentID: "<msg-txt>", MediaType: "txt"}}

	for idx, image := range images {
		mediaHeader := textproto.MIMEHeader{}

		contentType := http.DetectContentType(image)
		ext := getExtensionByContentType(http.DetectContentType(image))

		mediaHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s-%d\".", "image", idx))
		mediaHeader.Set("Content-ID", fmt.Sprintf("%s-%d%s", "image", idx, ext))
		mediaHeader.Set("Content-Type", contentType)

		attachments = append(attachments, SMILMedia{ContentID: fmt.Sprintf("<%s-%d>", "image", idx), MediaType: "img"})
		mediaPart, err := content.CreatePart(mediaHeader)
		if err != nil {
			return nil, "", err
		}

		_, err = io.Copy(mediaPart, bytes.NewReader(image))
		if err != nil {
			return nil, "", err
		}
	}

	// SMIL part.
	smildata, err := renderSMIL(attachments)
	if err != nil {
		return nil, "", err
	}

	smilHeader := textproto.MIMEHeader{}
	smilHeader.Set("Content-Type", "application/smil")
	smilHeader.Set("Content-ID", "<mms.smil>")
	smilPart, err := content.CreatePart(smilHeader)
	if err != nil {
		return nil, "", err
	}
	_, err = smilPart.Write(smildata)
	if err != nil {
		return nil, "", err
	}

	metaDataContent := textproto.MIMEHeader{}
	metaDataContent.Set("Content-Type", fmt.Sprintf("multipart/related; boundary=\"%s\"; start=\"%s\"; type=\"text/xml\"", content.Boundary(), "<mms.siml>"))
	metaDataContent.Set("Content-ID", "<mmscontent>")
	err = content.Close()
	if err != nil {
		return nil, "", err
	}

	contentPart, err := main.CreatePart(metaDataContent)
	if err != nil {
		return nil, "", err
	}
	_, err = contentPart.Write(contentBody.Bytes())
	if err != nil {
		return nil, "", err
	}

	// Close multipart writer.
	err = main.Close()
	if err != nil {
		return nil, "", err
	}

	// Request Content-Type with boundary parameter.
	contentType := fmt.Sprintf("multipart/related; boundary=%s", main.Boundary())

	return body, contentType, nil
}

type DeliverResponseParams struct {
	TransactionID string
	StatusCode    string
	StatusText    string
}

func GenerateMM7DeliverResponse(params DeliverResponseParams, soaptmpl *template.Template) (*bytes.Buffer, string, error) {
	soapdata := &bytes.Buffer{}

	contentType := "application/xml; charset=utf-8"

	_ = soaptmpl.Execute(soapdata, params)

	return soapdata, contentType, nil
}

var mimeTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
}

func getExtensionByContentType(contentType string) string {
	return mimeTypes[strings.ToLower(contentType)]
}
