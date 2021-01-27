package api

import (
	"net/http"

	"github.com/burstsms/mtmo-tp/backend/lib/errorlib"
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
		switch err.Error() {
		case errorlib.ErrInvalidMobileNumber.Error(),
			errorlib.ErrInvalidPhoneNumber.Error(),
			errorlib.ErrInvalidSenderNotFound.Error(),
			errorlib.ErrInvalidSenderChannel.Error(),
			errorlib.ErrInvalidSenderCountry.Error(),
			errorlib.ErrInvalidSMSTooManyParts.Error(),
			errorlib.ErrInsufficientBalance.Error(),
			errorlib.ErrInvalidRecipientInternationalNumber.Error():

			r.WriteError(err.Error(), http.StatusBadRequest)
		default:
			r.WriteError(err.Error(), http.StatusInternalServerError)
		}
		return
	}

	r.Write(res.SMS, http.StatusOK)
}
