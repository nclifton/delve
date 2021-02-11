package adminapi

import "net/http"

func ImportSenderPOST(r *Route) {
	type payload struct {
		Status string `json:"status"`
	}
	data := payload{
		Status: "ok",
	}

	r.Write(data, http.StatusOK)
}
