package mm7utils

import (
	"bytes"
	"text/template"
)

type SubmitResponseParams struct {
	TransactionID string
	StatusCode    string
	StatusText    string
	MessageID     string
}

func GenerateMM7SubmitResponse(params SubmitResponseParams, soaptmpl *template.Template) (*bytes.Buffer, string, error) {

	soapdata := &bytes.Buffer{}

	_ = soaptmpl.Execute(soapdata, params)

	// Request Content-Type with boundary parameter.
	contentType := "application/xml; charset=utf-8"

	return soapdata, contentType, nil
}
