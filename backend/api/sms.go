package api

import (
	"log"
	"net/http"
)

type SMSPOSTRequest struct {
	Message   string `json:"message" valid:"required"`
	Recipient string `json:"recipient"`
	Sender    string `json:"sender"`
}

func SMSPOST(r *Route) {
	account, err := r.RequireAccountContext()
	if err != nil {
		return
	}
	log.Printf("Account: %+v", account)

	var req SMSPOSTRequest
	err = r.DecodeRequest(&req)
	if err != nil {
		return
	}

	type payload struct {
		Message string `json:"message"`
	}

	data := payload{
		Message: req.Message,
	}

	r.Write(data, http.StatusOK)
}
