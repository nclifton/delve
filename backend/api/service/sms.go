package service

import (
	"net/http"

	"github.com/burstsms/mtmo-tp/backend/lib/errorlib"
	"github.com/burstsms/mtmo-tp/backend/lib/rest"
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

func (s *Service) SMSPOST(hc *rest.HandlerContext) {
	account := accountFromCtx(hc)

	var req SMSPOSTRequest
	if err := hc.DecodeJSON(&req); err != nil {
		return
	}

	res, err := s.SMSClient.Send(sms.SendParams{
		MessageRef: req.MessageRef,
		Message:    req.Message,
		AccountID:  account.GetId(),
		Sender:     req.Sender,
		Recipient:  req.Recipient,
		Country:    req.Country,
		AlarisUser: account.GetAlarisUsername(),
		AlarisPass: account.GetAlarisPassword(),
		AlarisURL:  account.GetAlarisUrl(),
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
			hc.WriteJSONError(err.Error(), http.StatusBadRequest, err)
		default:
			hc.LogFatal(err)
		}
		return
	}

	hc.WriteJSON(res.SMS, http.StatusOK)
}
