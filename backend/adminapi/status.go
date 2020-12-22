package adminapi

import "net/http"

func StatusGET(r *Route) {
	type payload struct {
		Status string `json:"status"`
	}
	data := payload{
		Status: "ok",
	}

	r.Write(data, http.StatusOK)
}
