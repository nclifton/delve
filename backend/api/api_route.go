package api

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	account "github.com/burstsms/mtmo-tp/backend/account/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/api/middleware/context"
	"github.com/burstsms/mtmo-tp/backend/lib/valid"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Route is some bullshit type so we can make httprouter chain with alice
// and still access params and such
type Route struct {
	params   httprouter.Params
	w        http.ResponseWriter
	r        *http.Request
	endpoint func(*Route)
	api      *API
}

// NewRoute returns a handler wrapped in middleware
// it's super important that route objects are created when requests come in
// if they are created before, the same object is shared for all requests
// and you fuck yourself with concurrent access... derp!
func NewRoute(api *API, chain alice.Chain, endpoint func(*Route)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		f := &Route{
			params:   p,
			endpoint: endpoint,
			api:      api,
		}

		chain.Then(f).ServeHTTP(w, r)
	}
}

// NewPlainRoute for plain http.Handler's
func NewPlainRoute(api *API, chain alice.Chain, endpoint func(*Route)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f := &Route{
			endpoint: endpoint,
			api:      api,
		}

		chain.Then(f).ServeHTTP(w, r)
	}
}

// EmptyRoute for convenience
func EmptyRoute(r *Route) {}

// NotFoundRoute route for convenience
func NotFoundRoute(r *Route) {
	r.WriteError("Not Found", http.StatusNotFound)
}

// IndexRoute route for convenience
func IndexRoute(r *Route) {
	r.w.Header().Set("X-Sendsei-Build", r.api.opts.Gitref)
	fmt.Fprint(r.w, "here be dragons\n")
}

// ServeHTTP for accepting wrapped objects
func (r *Route) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.w = w
	r.r = req
	r.endpoint(r)
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

	if err = valid.Validate(v); err != nil {
		log.Println(err)
		r.WriteValidatorError(err)
		return errors.New("request was invalid")
	}

	return nil
}

// Write sends the data in a json format if valid
// will also gzip encode the output if requested
func (r *Route) Write(v interface{}, code int) {
	var err error
	r.w.Header().Set("Content-Type", "application/json")

	if strings.Contains(r.r.Header.Get("Accept-Encoding"), "gzip") {
		gz := gzip.NewWriter(r.w)
		defer gz.Close()

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

// SuccessResponse type for a generic OK response
type SuccessResponse struct {
	Message string `json:"message"`
}

// WriteOK helper
func (r *Route) WriteOK() {
	r.Write(&SuccessResponse{"success"}, http.StatusOK)
}

// WriteText just for sending some text as a response
func (r *Route) WriteText(text string, code int) {
	r.w.WriteHeader(code)
	fmt.Fprintln(r.w, text)
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
	r.Write(&JSONErrors{Error: "Validation Error", ErrorData: valid.ErrorsByField(err)}, http.StatusOK)
}

// RequireAccountIDContext ensures we have authed on this request and account is loaded
// Returns the user and accountid
func (r *Route) RequireAccountContext() (*account.Account, error) {
	account, ok := context.Get(r.r, "account").(*account.Account)
	if !ok {
		return nil, errors.New("no such account")
	}

	return account, nil
}

// helper to fetch value from query string
func (r *Route) QueryParam(key string) string {
	return r.r.URL.Query().Get(key)
}

// helper to fetch integer from query string
func (r *Route) QueryParamInt(key string) (int, error) {
	val, err := strconv.ParseInt(r.r.URL.Query().Get(key), 10, 0)
	return int(val), err
}
