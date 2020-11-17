package api

import (
	"net/http"
)

type SMSPOSTRequest struct {
	Message   string `json:"message" valid:"required"`
	Recipient string `json:"recipient"`
	Sender    string `json:"sender"`
}

const (
	smsStatusSent  = "ok"
	smsStatusError = "error"
)

func SMSPOST(r *Route) {
	_, err := r.RequireAccountContext()
	if err != nil {
		return
	}

	var req SMSPOSTRequest
	err = r.DecodeRequest(&req)
	if err != nil {
		return
	}

	type payload struct {
		Message string
	}

	data := payload{
		Message: req.Message,
	}

	r.Write(data, http.StatusOK)
}
