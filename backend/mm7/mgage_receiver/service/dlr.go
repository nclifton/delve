package service

import (
	"log"
	"net/http"

	"github.com/burstsms/mtmo-tp/backend/lib/mgage"
	"github.com/burstsms/mtmo-tp/backend/lib/rest"
)

func (s *Service) DLRPOST(hc *rest.HandlerContext) {
	var req mgage.MTPOSTRequest
	if err := hc.DecodeJSON(&req); err != nil {
		return
	}

	// TODO: publish jobs to mgage receiver worker

	// TODO: validation on request?

	// TODO: do something with the request + return error as status ?
	log.Printf("RESULT: %+v", req)

	res := mgage.MTPOSTResponse{
		Status: "success",
	}

	hc.WriteJSON(res, http.StatusOK)
}
