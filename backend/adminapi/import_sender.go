package adminapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vincent-petithory/dataurl"
)

/**
 * api end point for parsing, validating and inserting sender data provided as a base64 encoded CSV data string
 */
func ImportSenderPOST(r *Route) {

	senders, failed := GetSendersFromRequest(r)
	if failed {
		return
	}

	for _, sender := range senders {
		log.Printf("%+v\n", sender)

		// use sender RPC service to insert the senders

	}

	type payload struct {
		Status string `json:"status"`
	}
	data := payload{
		Status: "ok",
	}

	r.Write(data, http.StatusOK)
}

func GetSendersFromRequest(r *Route) ([]SenderCSV, bool) {
	csvBytes, failed := getCSVFromRequest(r)
	if failed {
		return nil, true
	}

	senders := []SenderCSV{}

	err := gocsv.UnmarshalBytes(csvBytes, &senders)
	if err != nil {
		log.Println(err)
		r.WriteError("Failed to unmarshal the sender CSV: "+err.Error(), http.StatusBadRequest)
		return nil, true
	}
	return senders, false
}

type ImportJSON struct {
	Data string `json:"data"`
}
type SenderCSV struct {
	AccountId      string `csv:"account_id"`
	Address        string `csv:"address"`
	Country        string `csv:"country"`
	Channels       Array  `csv:"channels"`
	MMSProviderKey string `csv:"mms_provider_key"`
	Comment        string `csv:"comment"`
}

func getCSVFromRequest(r *Route) ([]byte, bool) {
	var j ImportJSON

	err := json.NewDecoder(r.r.Body).Decode(&j)
	if err != nil {
		log.Println(err)
		r.WriteError("Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return nil, true
	}

	csv, err := dataurl.DecodeString(j.Data)
	if err != nil {
		log.Println(err)
		r.WriteError("Invalid Data URL: "+err.Error(), http.StatusBadRequest)
		return nil, true
	}

	return csv.Data, false
}

type Array []string

// Convert the internal string array to JSON string
func (a *Array) MarshalCSV() (string, error) {
	str, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(str), nil
}

// Convert the CSV JSON string to string array (adminapi.Array)
func (a *Array) UnmarshalCSV(csv string) error {
	err := json.Unmarshal([]byte(csv), &a)
	return err
}

func (a *Array) String() []string {
	array := make([]string, len(*a))
	for _, str := range *a {
		array = append(array, str)
	}
	return array
}
