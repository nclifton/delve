package service

import (
	"net/http"

	"github.com/burstsms/mtmo-tp/backend/lib/errorlib"
	"github.com/burstsms/mtmo-tp/backend/lib/rest"
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

func (s *Service) MMSPOST(hc *rest.HandlerContext) {
	account := accountFromCtx(hc)

	var req MMSPOSTRequest
	if err := hc.DecodeJSON(&req); err != nil {
		return
	}

	reply, err := s.MMSClient.Send(mms.SendParams{
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
		// TODO: implement more errors for mms
		if err != nil {
			switch err.Error() {
			case errorlib.ErrInvalidSenderNotFound.Error(),
				errorlib.ErrInvalidSenderChannel.Error(),
				errorlib.ErrInvalidSenderCountry.Error(),
				errorlib.ErrInvalidSenderMMSProviderKeyEmpty.Error(),
				errorlib.ErrInvalidMMSLengthMessage.Error(),
				errorlib.ErrInvalidMMSLengthContentURLs.Error(),
				errorlib.ErrInvalidRecipientInternationalNumber.Error():
				hc.WriteJSONError(err.Error(), http.StatusBadRequest, err)
			default:
				hc.LogFatal(err)
			}
			return
		}
	}

	hc.WriteJSON(reply.MMS, http.StatusOK)
}
