package tecloo

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/mm7utils"
	"github.com/google/uuid"
)

var (
	stripRegex         = regexp.MustCompile(`>\s+<`)
	submitRegex        = regexp.MustCompile(`<SubmitReq.*>`)
	recipientRegex     = regexp.MustCompile(`<Recipients><To><Number>(.*)<\/Number><\/To></Recipients>`)
	senderRegex        = regexp.MustCompile(`<SenderAddress><ShortCode>(.*)<\/ShortCode><\/SenderAddress>`)
	transactionIDRegex = regexp.MustCompile(`<TransactionID.*>(.*)<\/TransactionID>`)
	subjectRegex       = regexp.MustCompile(`<Subject>(.*)<\/Subject>`)
)

func SubmitPOST(r *Route) {
	contentType, response, err := r.api.parseSubmit(r.r.Header.Get(`Content-Type`), r.r.Body)
	if err != nil {
		r.WriteError(err.Error(), http.StatusBadRequest)
	}

	r.w.Header().Set("Content-Type", contentType)
	r.w.WriteHeader(http.StatusOK)
	_, err = r.w.Write(response.Bytes())
	if err != nil {
		r.WriteError("Could not write response output", http.StatusBadRequest)
		return
	}

}

func (api *TeclooAPI) parseSubmit(contentType string, body io.Reader) (string, *bytes.Buffer, error) {
	ctx := context.Background()

	parts, err := mm7utils.ProcessMultiPart(contentType, body)
	if err != nil {
		api.log.Errorf("Invalid Request: %s", err)
	}

	var recipient string
	var sender string
	var transactionid string
	var subject string
	var message string
	status := "1000" // default success
	var fileSums []map[string]string
	var last4 string

	for _, part := range parts {
		fileHash := sha1.New()

		api.log.Fields(ctx, logger.Fields{"ContentID": part.ContentID}).Debug("Processing Part")

		if strings.Trim(part.ContentID, "<>") == "soap-start" {

			soap := strings.Replace(string(part.Body), "\n", "", -1)
			soap = stripRegex.ReplaceAllString(soap, "><")
			if !submitRegex.MatchString(soap) {
				return "", nil, errors.New("Submit Request must contain a SubmitReq Element")
			}
			recipient = mm7utils.ExtractEntity(*recipientRegex, soap)
			// Ok lets set the status based on the last 4 digits of the recipient phone number
			last4 = recipient[len(recipient)-4:]
			if _, ok := statusCodes[last4]; ok {
				status = last4
			}
			sender = mm7utils.ExtractEntity(*senderRegex, soap)
			subject = mm7utils.ExtractEntity(*subjectRegex, soap)
			transactionid = mm7utils.ExtractEntity(*transactionIDRegex, soap)
			continue
		}

		if strings.HasPrefix(part.ContentType, "image/") {
			_, err := fileHash.Write(part.Body)
			if err != nil {
				return "", nil, errors.New("Could not hash image content")
			}
			fileSum := hex.EncodeToString(fileHash.Sum(nil))
			fileSums = append(fileSums, map[string]string{part.ContentID: fileSum})
			continue
		}

		if strings.HasPrefix(part.ContentType, "text/plain") {
			message = string(part.Body)
			continue
		}

	}
	api.log.Fields(ctx, logger.Fields{
		"TransactionID": transactionid,
		"Subject":       subject,
		"Message":       message,
		"FileSum":       fileSums,
		"Recipient":     recipient,
		"Sender":        sender,
		"Status":        status,
		"Status Text":   statusCodes[status].text,
	}).Info("MM7 Submit")

	uuid := uuid.New()
	if err != nil {
		return "", nil, errors.New("Could not generate a message id")
	}

	response, contentType, _ := mm7utils.GenerateMM7SubmitResponse(mm7utils.SubmitResponseParams{
		StatusCode:    status,
		StatusText:    statusCodes[status].text,
		TransactionID: transactionid,
		MessageID:     uuid.String()}, api.templates.SubmitResponse)

	if status == "1000" {
		api.sendDRRequest(&DRParams{
			TransactionID: transactionid,
			Recipient:     recipient,
			Sender:        sender,
			Status:        last4,
			MessageID:     uuid.String(),
		})
	}

	return contentType, response, nil
}
