package rest

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/burstsms/mtmo-tp/backend/lib/valid"
)

type Validator = func(v interface{}, cvs ...valid.CustomValidator) error

// HandlerContext injected into handlers/routes, provides various utils for I/O and logging
type HandlerContext struct {
	w    http.ResponseWriter
	r    *http.Request
	opts *HandlerConfig
}

func NewHandlerContext(w http.ResponseWriter, r *http.Request, opts *HandlerConfig) *HandlerContext {
	return &HandlerContext{w: w, r: r, opts: opts}
}

func (hc *HandlerContext) LogError(err error) {
	hc.opts.Log.Error(hc.r.Context(), fmt.Sprintf("%s %s", hc.r.Method, hc.r.RequestURI), err.Error())
}

func (hc *HandlerContext) LogInfo(msg string) {
	hc.opts.Log.Info(hc.r.Context(), fmt.Sprintf("%s %s", hc.r.Method, hc.r.RequestURI), msg)
}

func (hc *HandlerContext) LogFatal(err error) {
	hc.LogError(err)
	panic(err)
}

func (hc *HandlerContext) Context() context.Context {
	return hc.r.Context()
}

func (hc *HandlerContext) FromContext(key interface{}) interface{} {
	return hc.r.Context().Value(key)
}

func (hc *HandlerContext) Params() httprouter.Params {
	return hc.FromContext(ParamsKey).(httprouter.Params)
}

func (hc *HandlerContext) DecodeJSON(v interface{}) error {
	if hc.r.Header.Get("Content-Type") != "application/json" {
		err := errors.New("Content-Type must be application/json")
		hc.WriteJSONError(err.Error(), http.StatusBadRequest, err)
		return err
	}

	err := json.NewDecoder(hc.r.Body).Decode(v)
	if err != nil {
		hc.WriteJSONError(fmt.Sprintf("Invalid JSON: %s", err.Error()), http.StatusBadRequest, err)
		return err
	}

	if hc.opts.JSONValidator != nil {
		if err := hc.opts.JSONValidator(v); err != nil {
			hc.WriteJSONError(err.Error(), http.StatusBadRequest, err)
			return err
		}
	}

	return nil
}

func (hc *HandlerContext) WriteJSON(v interface{}, code int) {
	var err error
	hc.w.Header().Set("Content-Type", "application/json")

	if strings.Contains(hc.r.Header.Get("Accept-Encoding"), "gzip") {
		gz := gzip.NewWriter(hc.w)
		defer gz.Close()

		hc.w.Header().Set("Content-Encoding", "gzip")
		hc.w.WriteHeader(code)
		err = json.NewEncoder(gz).Encode(v)
	} else {
		hc.w.WriteHeader(code)
		err = json.NewEncoder(hc.w).Encode(v)
	}

	if err != nil {
		hc.LogFatal(err)
	}
}

func (hc *HandlerContext) WriteJSONError(msg string, code int, err error) {
	type jsonError struct {
		Error string `json:"error"`
	}

	if err != nil {
		hc.LogError(err)
	}

	hc.WriteJSON(&jsonError{Error: msg}, code)
}

func (hc *HandlerContext) WriteJSONSuccess(msg string) {
	type jsonSuccess struct {
		Message string `json:"message"`
	}

	hc.WriteJSON(&jsonSuccess{Message: msg}, 200)
}
