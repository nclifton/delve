package tecloo_receiver

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	//"github.com/burstsms/mtmo-tp/backend/biz/valid"
	"github.com/burstsms/mtmo-tp/backend/lib/mongo"
)

// Write sends the data in a json format if valid
// will also gzip encode the output if requested
func (r *Route) Write(v interface{}, code int) {
	var err error
	r.w.Header().Set("Content-Type", "application/json")

	if strings.Contains(r.r.Header.Get("Accept-Encoding"), "gzip") {
		gz := gzip.NewWriter(r.w)
		defer func() {
			err := gz.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		r.w.Header().Set("Content-Encoding", "gzip")
		r.w.WriteHeader(code)
		err = json.NewEncoder(gz).Encode(v)
	} else {
		r.w.WriteHeader(code)
		err = json.NewEncoder(r.w).Encode(v)
	}

	if err != nil {
		log.Println("ERROR: ENCODING JSON:", err)
	}
}

// DecodeRequest helper to parse and validate json body in request
func (r *Route) DecodeWithoutValidatingRequest(v interface{}) error {
	if r.r.Header.Get("Content-Type") != "application/json" {
		r.WriteError("Content-Type must be application/json", http.StatusBadRequest)
		return errors.New("Expected json content type")
	}

	err := json.NewDecoder(r.r.Body).Decode(v)
	if err != nil {
		log.Println(err)
		r.WriteError("Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return err
	}

	return nil
}

// DecodeAndValidateRequest will make sure the json input conforms to validators
func (r *Route) DecodeRequest(v interface{}) error {
	err := r.DecodeWithoutValidatingRequest(v)
	if err != nil {
		return err
	}

	// if err = valid.Validate(v); err != nil {
	// 	log.Println(err)
	// 	r.WriteValidatorError(err)
	// 	return errors.New("request was invalid")
	// }

	return nil
}

// JSONError for returning error strings in json payload
type JSONError struct {
	Error string `json:"error"`
}

// WriteError helper for just sending an error string
func (r *Route) WriteError(err string, code int) {
	r.Write(&JSONError{Error: err}, code)
}

type JSONErrors struct {
	Error     string            `json:"error"`
	ErrorData map[string]string `json:"error_data"`
}

func (r *Route) WriteValidatorError(err error) {
	//r.Write(&JSONErrors{Error: "Validation Error", ErrorData: valid.ErrorsByField(err)}, http.StatusOK)
}

// for checking id and writing errors out
func (r *Route) EnsureIDParam() (mongo.OID, error) {
	return r.EnsureIDParamByName("id")
}

func (r *Route) EnsureIDParamByName(name string) (mongo.OID, error) {
	idStr := r.params.ByName(name)
	id, err := mongo.OIDFromHex(idStr)
	if err != nil {
		r.WriteError("Invalid ID: "+idStr, http.StatusBadRequest)
		return id, errors.New("Invalid ID")
	}

	return id, nil
}
