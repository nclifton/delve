package tecloo_receiver

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/burstsms/mtmo-tp/backend/lib/mm7utils"
)

var (
	stripRegex                      = regexp.MustCompile(`>\s+<`)
	deliverRegex                    = regexp.MustCompile(`<DeliverReq.*>`)
	deliveryReportRegex             = regexp.MustCompile(`<DeliveryReportReq.*>`)
	deliverReqRecipientRegex        = regexp.MustCompile(`<Recipients.*><To.*><Number.*>(.*)<\/Number><\/To></Recipients>`)
	deliveryReportReqRecipientRegex = regexp.MustCompile(`<Recipient.*><Number.*>(.*)<\/Number><\/Recipient>`)
	senderRegex                     = regexp.MustCompile(`<Sender.*><Number.*>(.*)<\/Number></Sender>`)
	transactionIDRegex              = regexp.MustCompile(`<TransactionID.*>(.*)<\/TransactionID>`)
	linkedIDRegex                   = regexp.MustCompile(`<LinkedID.*>(.*)<\/LinkedID>`)
	timestampRegex                  = regexp.MustCompile(`<TimeStamp.*>(.*)<\/TimeStamp>`)
	statusCodeRegex                 = regexp.MustCompile(`<StatusCode.*>(.*)<\/StatusCode>`)
	statusTextRegex                 = regexp.MustCompile(`<StatusText.*>(.*)<\/StatusText>`)
	messageIDRegex                  = regexp.MustCompile(`<MessageID.*>(.*)<\/MessageID>`)
	dateRegex                       = regexp.MustCompile(`<Date.*>(.*)<\/Date>`)
	mmstatusRegex                   = regexp.MustCompile(`<MMStatus.*>(.*)<\/MMStatus>`)
)

func extractEntity(regex regexp.Regexp, soap string) string {
	matches := regex.FindStringSubmatch(soap)
	if len(matches) > 0 {
		return matches[1]
	}
	return ""
}

// receives external soap request in mm7 format
// validates the request and puts it into inbound job
// return relevant http response to tecloo
func InboundPOST(r *Route) {
	var reqBytes bytes.Buffer
	status := "1000"
	body := io.TeeReader(r.r.Body, &reqBytes)

	requestContentType := r.r.Header.Get("Content-Type")

	// if it's a delivery request
	if strings.HasPrefix(requestContentType, "multipart/") {
		status, err := processDeliverRequest(requestContentType, body)
		if err != nil {
			r.api.log.Errorf("Invalid Request: %s", err)
		}

		response, contentType, _ := mm7utils.GenerateMM7DeliverResponse(mm7utils.DeliverResponseParams{
			StatusCode: status,
			StatusText: mm7utils.StatusCodes[status].Text,
		}, r.api.templates.DeliverResponse)

		r.w.Header().Set("Content-Type", contentType)
		r.w.WriteHeader(http.StatusOK)
		_, err = r.w.Write(response.Bytes())
		if err != nil {
			r.WriteError("Could not write response output", http.StatusBadRequest)
			return
		}
	}

	// if it's a delivery report request
	if strings.HasPrefix(requestContentType, "application/") {
		status, transactionID, err := processDeliveryReportRequest(requestContentType, body)
		if err != nil {
			r.api.log.Errorf("Invalid Request: %s", err)
		}

		response, contentType, _ := mm7utils.GenerateMM7DeliveryReportResponse(mm7utils.DeliveryReportResponseParams{
			TransactionID: transactionID,
			StatusCode:    status,
			StatusText:    mm7utils.StatusCodes[status].Text,
		}, r.api.templates.DeliveryReportResponse)

		r.w.Header().Set("Content-Type", contentType)
		r.w.WriteHeader(http.StatusOK)
		_, err = r.w.Write(response.Bytes())
		if err != nil {
			r.WriteError("Could not write response output", http.StatusBadRequest)
			return
		}
	}

	if status == "1000" {
		// push the full multipart request in a job onto queue for dl/dlr worker
		_ = r.api.Publish(reqBytes.Bytes(), map[string]interface{}{"Content-Type": requestContentType}, true, "dldr.fake")
	}
}

func processDeliverRequest(contentType string, body io.Reader) (statusCode string, err error) {
	parts, err := mm7utils.ProcessMultiPart(contentType, body)
	if err != nil {
		return "2007", err
	}

	var recipient string
	status := "1000" // default success
	var last4 string

	for _, part := range parts {
		if part.ContentID == "<soap-start>" {
			soap := strings.Replace(string(part.Body), "\n", "", -1)
			soap = stripRegex.ReplaceAllString(soap, "><")

			if !deliverRegex.MatchString(soap) {
				return "2007", errors.New("Request not recognised")
			}

			if !linkedIDRegex.MatchString(soap) {
				return "2007", errors.New("Deliver Request must contain a LinkedID Element")
			}

			if !timestampRegex.MatchString(soap) {
				return "2007", errors.New("Deliver Request must contain a TimeStamp Element")
			}

			if !transactionIDRegex.MatchString(soap) {
				return "2007", errors.New("Request must contain a TransactionID Element")
			}

			if !deliverReqRecipientRegex.MatchString(soap) {
				return "2007", errors.New("Request must contain a Recipients Element")
			}
			recipient = extractEntity(*deliverReqRecipientRegex, soap)

			if !senderRegex.MatchString(soap) {
				return "2007", errors.New("Request must contain a Sender Element")
			}

			// Ok lets set the status based on the last 4 digits of the recipient phone number
			last4 = recipient[len(recipient)-4:]
			if _, ok := mm7utils.StatusCodes[last4]; ok {
				status = last4
			}
			continue
		}
	}

	return status, nil
}

func processDeliveryReportRequest(contentType string, body io.Reader) (statusCode, transactionID string, err error) {
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return "2007", "", errors.New("Unable to read request bytes")
	}

	var recipient string
	status := "1000" // default success
	var last4 string

	soap := strings.Replace(string(bodyBytes), "\n", "", -1)
	soap = stripRegex.ReplaceAllString(soap, "><")

	if !deliveryReportRegex.MatchString(soap) {
		return "2007", "", errors.New("Request not recognised")
	}

	if !messageIDRegex.MatchString(soap) {
		return "2007", "", errors.New("Delivery Report Request must contain a message ID Element")
	}

	if !dateRegex.MatchString(soap) {
		return "2007", "", errors.New("Delivery Report Request must contain a date Element")
	}

	if !mmstatusRegex.MatchString(soap) {
		return "2007", "", errors.New("Delivery Report Request must contain a MMStatus Element")
	}

	if !transactionIDRegex.MatchString(soap) {
		return "2007", "", errors.New("Request must contain a TransactionID Element")
	}
	transactionID = extractEntity(*transactionIDRegex, soap)

	if !deliveryReportReqRecipientRegex.MatchString(soap) {
		return "2007", "", errors.New("Request must contain a Recipients Element")
	}
	recipient = extractEntity(*deliveryReportReqRecipientRegex, soap)

	if !senderRegex.MatchString(soap) {
		return "2007", "", errors.New("Request must contain a Sender Element")
	}

	// Ok lets set the status based on the last 4 digits of the recipient phone number
	last4 = recipient[len(recipient)-4:]
	if _, ok := mm7utils.StatusCodes[last4]; ok {
		status = last4
	}

	return status, transactionID, nil
}
