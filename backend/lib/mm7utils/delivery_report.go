package mm7utils

import (
	"bytes"
	"text/template"
)

type DeliveryReportParams struct {
	TransactionID string
	StatusCode    string
	StatusText    string
	Date          string
	MessageID     string
	Recipient     string
	Sender        string
}

func GenerateMM7DeliveryReport(params DeliveryReportParams, soaptmpl *template.Template) (*bytes.Buffer, string, error) {

	soapdata := &bytes.Buffer{}

	if err := soaptmpl.Execute(soapdata, params); err != nil {
		return nil, "", err
	}

	// Request Content-Type with boundary parameter.
	contentType := "application/xml; charset=utf-8"

	return soapdata, contentType, nil
}

type DeliveryReportResponseParams struct {
	TransactionID string
	StatusCode    string
	StatusText    string
}

func GenerateMM7DeliveryReportResponse(params DeliveryReportResponseParams, soaptmpl *template.Template) (*bytes.Buffer, string, error) {
	soapdata := &bytes.Buffer{}

	_ = soaptmpl.Execute(soapdata, params)

	// Request Content-Type with boundary parameter.
	contentType := "text/xml; charset=utf-8"

	return soapdata, contentType, nil
}
