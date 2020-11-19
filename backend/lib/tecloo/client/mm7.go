package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/burstsms/mtmo-tp/backend/lib/mm7utils"
)

type PostMM7Params struct {
	ID        string
	Subject   string
	Message   string
	Sender    string
	Recipient string
	Images    [][]byte
}

type PostMM7Response struct {
	XMLName xml.Name  `xml:"Envelope"`
	Text    string    `xml:",chardata"`
	SoapEnv string    `xml:"soap-env,attr"`
	Xmlns   string    `xml:"xmlns,attr"`
	Header  MM7Header `xml:"Header"`
	Body    MM7Body   `xml:"Body"`
}

type MM7Header struct {
	Text          string        `xml:",chardata"`
	TransactionID TransactionID `xml:"TransactionID"`
}

type TransactionID struct {
	Text           string `xml:",chardata"`
	MustUnderstand string `xml:"mustUnderstand,attr"`
}

type MM7Body struct {
	Text      string    `xml:",chardata"`
	SubmitRsp SubmitRsp `xml:"SubmitRsp"`
}

type SubmitRsp struct {
	Text       string `xml:",chardata"`
	MM7Version string `xml:"MM7Version"`
	Status     Status `xml:"Status"`
	MessageID  string `xml:"MessageID"`
}

type Status struct {
	Text       string `xml:",chardata"`
	StatusCode string `xml:"StatusCode"`
	StatusText string `xml:"StatusText"`
}

func (s service) PostMM7(params PostMM7Params, soaptmpl *template.Template) (PostMM7Response, int, error) {
	submit := mm7utils.SubmitParams{
		TransactionID:    params.ID,
		Subject:          params.Subject,
		VASPID:           "MYVASPID",
		Sender:           params.Sender,
		Recipient:        params.Recipient,
		AllowAdaptations: true,
	}

	soapBody, contentType, err := mm7utils.GenerateMM7Submit(submit, soaptmpl, params.Message, params.Images)
	if err != nil {
		return PostMM7Response{}, 0, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/mm7", s.apiURL), bytes.NewReader(soapBody.Bytes()))
	if err != nil {
		return PostMM7Response{}, 0, err
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := client.Do(req)
	if err != nil {
		return PostMM7Response{}, 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return PostMM7Response{}, resp.StatusCode, fmt.Errorf("Unexpected status: %s", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return PostMM7Response{}, resp.StatusCode, err
	}

	result := PostMM7Response{}
	if err = xml.Unmarshal(b, &result); err != nil {
		return PostMM7Response{}, resp.StatusCode, fmt.Errorf("%s while unmarshalling response body: %s", err, string(b))
	}

	return result, resp.StatusCode, nil
}
