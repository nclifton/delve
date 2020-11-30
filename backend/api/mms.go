package api

import (
	"fmt"
	"net/http"

	"github.com/burstsms/mtmo-tp/backend/mms/rpc"
	mms "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
)

func MMSGET(r *Route) {
	account, err := r.RequireAccountContext()
	if err != nil {
		return
	}

	id := r.params.ByName("id")

	mms, err := r.api.mms.FindByID(id, account.ID)
	if err != nil {
		r.WriteError("failed fetching MMS", http.StatusInternalServerError)
		return
	}

	type payload struct {
		// TODO: expose only the client method?
		MMS *rpc.MMS `json:"mms"`
	}

	data := payload{
		MMS: mms.MMS,
	}

	r.Write(data, http.StatusOK)
}

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

	validSender := checkValidSender(req.Sender, account.SenderMMS)
	if !validSender {
		r.WriteError(fmt.Sprintf("Sender: %s is not valid for account: %s(%s)", req.Sender, account.Name, account.ID), http.StatusBadRequest)
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
	if err != nil {
		r.WriteError(fmt.Sprintf("failed sending MMS: %s", err), http.StatusInternalServerError)
		return
	}

	r.Write(res.MMS, http.StatusOK)
}
