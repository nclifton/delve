package api

import (
	"fmt"
	"net/http"

	"github.com/burstsms/mtmo-tp/backend/lib/errorlib"
	mms "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
)

type MMSPOSTRequest struct {
	Subject     string   `json:"subject"`
	Message     string   `json:"message"`
	Recipient   string   `json:"recipient"`
	Sender      string   `json:"sender"`
	Country     string   `json:"country"`
	MessageRef  string   `json:"message_ref"`
	ContentURLs []string `json:"content_urls"`
	TrackLinks  bool     `json:"track_links"`
}

func MMSPOST(r *Route) {
	account, err := r.RequireAccountContext()
	if err != nil {
		r.WriteError(fmt.Sprintf("Could not process mms: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var req MMSPOSTRequest
	if err := r.DecodeRequest(&req); err != nil {
		r.WriteError(fmt.Sprintf("Could not process mms: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	providerKey := account.MMSProviderKey
	if providerKey == "" {
		r.WriteError("failed sending MMS Incorrectly configured provider", http.StatusInternalServerError)
		return
	}

	res, err := r.api.mms.Send(mms.SendParams{
		AccountID:   account.ID,
		Subject:     req.Subject,
		Message:     req.Message,
		Recipient:   req.Recipient,
		Sender:      req.Sender,
		Country:     req.Country,
		MessageRef:  req.MessageRef,
		ContentURLs: req.ContentURLs,
		TrackLinks:  req.TrackLinks,
	})

	// TODO: implement more errors for mms
	if err != nil {
		switch err {
		case errorlib.ErrInvalidSenderNotFound,
			errorlib.ErrInvalidSenderChannel,
			errorlib.ErrInvalidSenderCountry:
			r.WriteError(err.Error(), http.StatusBadRequest)
		default:
			r.WriteError(err.Error(), http.StatusInternalServerError)
		}
		return
	}

	r.Write(res.MMS, http.StatusOK)
}
