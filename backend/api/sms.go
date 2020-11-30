package api

import (
	"fmt"
	"net/http"

	sms "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
)

type SMSPOSTRequest struct {
	Message    string `json:"message" valid:"required"`
	MessageRef string `json:"message_ref"`
	Recipient  string `json:"recipient" valid:"required"`
	Sender     string `json:"sender" valid:"required"`
	Country    string `json:"country"`
	TrackLinks bool   `json:"track_links"`
}

func SMSPOST(r *Route) {
	account, err := r.RequireAccountContext()
	if err != nil {
		return
	}

	var req SMSPOSTRequest
	err = r.DecodeRequest(&req)
	if err != nil {
		return
	}

	// Check sender against the acconts sender list
	var validSender bool
	for _, possibleSender := range account.SenderSMS {
		validSender = possibleSender == req.Sender
	}

	if !validSender {
		r.WriteError(fmt.Sprintf("Sender: %s is not valid for account: %s(%s)", req.Sender, account.Name, account.ID), http.StatusBadRequest)
		return
	}

	res, err := r.api.sms.Send(sms.SendParams{
		MessageRef: req.MessageRef,
		Message:    req.Message,
		AccountID:  account.ID,
		Sender:     req.Sender,
		Recipient:  req.Recipient,
		Country:    req.Country,
		AlarisUser: account.AlarisUsername,
		AlarisPass: account.AlarisPassword,
		AlarisURL:  account.AlarisURL,
		TrackLinks: req.TrackLinks,
	})
	if err != nil {
		// handler rpc error
		r.WriteError(fmt.Sprintf("Could not process sms: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	r.Write(res.SMS, http.StatusOK)
}
