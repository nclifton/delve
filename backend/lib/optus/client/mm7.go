package client

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"text/template"
	"time"

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
		VASPID:           "6001",
		VASID:            "6001",
		Sender:           params.Sender,
		Recipient:        params.Recipient,
		AllowAdaptations: true,
	}

	soapBody, contentType, err := mm7utils.GenerateMM7Submit(submit, soaptmpl, params.Message, params.Images)
	if err != nil {
		return PostMM7Response{}, 0, fmt.Errorf("Could not generate MM7 Submit body: %s", err)
	}

	client := &http.Client{
		Timeout: time.Second * 30,
	}
	log.Printf("Optus URL: %s", s.apiURL)
	req, err := http.NewRequest("POST", s.apiURL, bytes.NewReader(soapBody.Bytes()))
	if err != nil {
		return PostMM7Response{}, 0, fmt.Errorf("Could not create post request: %s", err)
	}

	basicAuthStr := fmt.Sprintf("Basic %s", s.auth)
	req.Header.Set("Content-Type", contentType)
	req.Header.Add("Authorization", basicAuthStr)
	//req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Length", fmt.Sprintf("%d", len(soapBody.Bytes())))
	req.Header.Add("SOAPAction", "")
	req.Header.Add("TE", "deflate,gzip;q=0.3")
	req.Header.Add("User-Agent", "Sendsei/1.0")
	req.Header.Add("Connection", "TE, close")

	// Save a copy of this request for debugging.
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))

	resp, err := client.Do(req)
	if err != nil {
		return PostMM7Response{}, 0, fmt.Errorf("Could not perform request: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return PostMM7Response{}, resp.StatusCode, fmt.Errorf("Unexpected status: %s", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return PostMM7Response{}, resp.StatusCode, fmt.Errorf("Could not read response body: %s", err)
	}

	log.Printf("Optus Response Body: %s", string(b))

	result := PostMM7Response{}
	if err = xml.Unmarshal(b, &result); err != nil {
		return PostMM7Response{}, resp.StatusCode, fmt.Errorf("%s while unmarshalling response body: %s", err, string(b))
	}

	return result, resp.StatusCode, nil
}
