package tecloo

import (
	"bytes"
	"net/http"
	"time"

	"github.com/burstsms/mtmo-tp/backend/lib/mm7utils"
	"github.com/burstsms/mtmo-tp/backend/logger"
)

type DRParams struct {
	TransactionID string
	Recipient     string
	Sender        string
	Status        string
	MessageID     string
}

var drStatusCodes = map[string]string{
	"1000": "Retrieved",
	"1001": "Deferred",
	"1002": "Expired",
	"1003": "Forwarded",
	"1004": "Indeterminate",
	"1005": "NotSupported",
	"1006": "Rejected",
	"1007": "Retrieved",
	"1008": "Unreachable",
	"1009": "Unrecognized",
}

func (api *TeclooAPI) sendDRRequest(params *DRParams) {

	var status = "1000"

	if _, ok := drStatusCodes[params.Status]; ok {
		status = params.Status
	}

	api.log.Fields(logger.Fields{
		"TransactionID": params.TransactionID,
		"Recipient":     params.Recipient,
		"Sender":        params.Sender,
		"Status":        status,
		"StatusText":    drStatusCodes[status],
		"MessageID":     params.MessageID,
	}).Info("MM7 DR")

	if api.opts.DREndpoint != "" {
		body, contentType, _ := mm7utils.GenerateMM7DeliveryReport(mm7utils.DeliveryReportParams{
			StatusCode:    status,
			StatusText:    drStatusCodes[status],
			TransactionID: params.TransactionID,
			Sender:        params.Sender,
			Recipient:     params.Recipient,
			Date:          time.Now().UTC().Format(time.RFC3339),
			MessageID:     params.MessageID}, api.templates.SendDeliveryReport)

		req, err := http.NewRequest("POST", api.opts.DREndpoint, bytes.NewReader(body.Bytes()))
		req.Header.Set("Content-Type", contentType)
		if err != nil {
			api.log.Errorf("Could not create DR request: %s", err)
		}

		_, err = api.client.Do(req)
		if err != nil {
			api.log.Errorf("Could not do DR request: %s", err)
		}
	}

}
