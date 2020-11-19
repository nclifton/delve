package tecloo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/burstsms/mtmo-tp/backend/lib/mm7utils"
	"github.com/burstsms/mtmo-tp/backend/logger"
)

type HandsetParams struct {
	ID          string
	Subject     string
	Message     string
	Sender      string
	Recipient   string
	ContentURLs []string
	ProviderKey string
}

func HandsetPOST(r *Route) {
	var req HandsetParams
	err := r.DecodeRequest(&req)
	if err != nil {
		return
	}

	delivery := mm7utils.DeliverRequestParams{
		TransactionID:    req.ID,
		Subject:          req.Subject,
		Sender:           req.Sender,
		Recipient:        req.Recipient,
		AllowAdaptations: true,
	}

	var images [][]byte

	for _, url := range req.ContentURLs {
		response, err := http.Get(url)
		if err != nil {
			r.WriteError(fmt.Sprintf("Could not get content from url: %s %s ", url, err.Error()), http.StatusBadRequest)
		}
		defer func() {
			err := response.Body.Close()
			if err != nil {
				r.WriteError(fmt.Sprintf("Could not close response body for content from url: %s %s ", url, err.Error()), http.StatusBadRequest)
			}
		}()

		imageData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			r.WriteError(fmt.Sprintf("Could not read content from url: %s %s ", url, err.Error()), http.StatusBadRequest)
			return
		}

		images = append(images, imageData)
	}

	// Send to dr endpoint
	if r.api.opts.DREndpoint != "" {
		body, contentType, err := mm7utils.GenerateMM7DeliverRequest(delivery, r.api.templates.SendDelivery, req.Message, images)
		if err != nil {
			r.api.log.Errorf("Could not generate MM7 Delivery body: %s", err)
		}

		req, err := http.NewRequest("POST", r.api.opts.DREndpoint, bytes.NewReader(body.Bytes()))
		req.Header.Set("Content-Type", contentType)
		if err != nil {
			r.api.log.Errorf("Could not create DR request: %s", err)
		}

		_, err = r.api.client.Do(req)
		if err != nil {
			r.api.log.Errorf("Could not do DR request: %s", err)
		}
	}
	r.api.log.Fields(logger.Fields{
		"TransactionID": req.ID,
		"Subject":       req.Subject,
		"Message":       req.Message,
		"Recipient":     req.Recipient,
		"Sender":        req.Sender,
		"Status":        "1000",
		"Status Text":   statusCodes["1000"].text,
	}).Info("MM7 Sent")

	type payload struct {
		Result bool `json:"result"`
	}

	data := payload{
		Result: true,
	}

	r.Write(data, http.StatusOK)

}
