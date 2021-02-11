package adminapi

import "net/http"

/**
 * api end point for parsing, validating and inserting sender data provided as a base64 encoded CSV data string
 */
func ImportSenderPOST(r *Route) {

	







	type payload struct {
		Status string `json:"status"`
	}
	data := payload{
		Status: "ok",
	}

	r.Write(data, http.StatusOK)
}
