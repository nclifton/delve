package adminapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/grpc/status"

	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
)

type ImportSenderPOSTRequest struct {
	Data []byte `json:"data"`
}

/**
 * api end point for parsing, validating and inserting sender data provided as a base64 encoded CSV data string
 */
func ImportSenderPOST(r *Route) {

	var importSenderPOSTRequest ImportSenderPOSTRequest
	err := json.NewDecoder(r.r.Body).Decode(&importSenderPOSTRequest)

	if err != nil {
		log.Println(err)
		r.WriteError("Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// discarding the return from the RPC call for now until it gets better defined
	_, err = r.api.sender.CreateSendersFromCSVDataURL(r.r.Context(),
		&senderpb.CreateSendersFromCSVDataURLParams{
			CSV: importSenderPOSTRequest.Data,
		})
	if err != nil {
		// handler rpc error
		grpcError := status.Convert(err)
		log.Printf("Could not upload senders CSV: %s", err.Error())
		r.WriteError(fmt.Sprintf("Could not upload senders CSV: %s", grpcError.Message()), http.StatusInternalServerError)
		return
	}

	type payload struct {
		Status string `json:"status"`
	}
	data := payload{
		Status: "ok",
	}
	r.Write(data, http.StatusOK)

}
