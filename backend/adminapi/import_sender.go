package adminapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/grpc/status"

	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
)

func ImportSenderPOST(r *Route) {

	senders, failed := getCSVFromRequestBodyJSON(r)
	if failed {
		return
	}

	// discarding the return from the RPC call for now until it gets better defined
	_, err := r.api.sender.CreateSendersFromCSVDataURL(r.r.Context(),
		&senderpb.CreateSendersFromCSVDataURLParams{
			CSV: senders,
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

type ImportJSON struct {
	Data []byte `json:"data"`
}

func getCSVFromRequestBodyJSON(r *Route) ([]byte, bool) {

	var j ImportJSON
	err := json.NewDecoder(r.r.Body).Decode(&j)
	if err != nil {
		log.Println(err)
		r.WriteError("Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return nil, true
	}

	return j.Data, false
}
