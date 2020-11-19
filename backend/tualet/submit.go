package tualet

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/burstsms/mtmo-tp/backend/logger"
)

type submitParams struct {
	username        string
	password        string
	command         string
	message         string
	dnis            string
	ani             string
	longMessageMode string
}

func checkParams(params url.Values) (int, string, submitParams) {

	status := http.StatusOK
	response := "submission"

	values := submitParams{}

	values.username = params.Get("username")
	if values.username == "" {
		status = http.StatusUnauthorized
		response = "not authorized (check login and password)"
		return status, response, values
	}

	values.password = params.Get("password")
	if values.password == "" {
		status = http.StatusUnauthorized
		response = "not authorized (check login and password)"
		return status, response, values
	}

	values.command = params.Get("command")
	if values.command != "submit" {
		status = http.StatusBadRequest
		response = "invalid command"
		return status, response, values
	}

	values.message = params.Get("message")
	if values.message == "" {
		status = http.StatusBadRequest
		response = "message missing"
		return status, response, values
	}

	values.dnis = params.Get("dnis")
	if values.dnis == "" {
		status = http.StatusBadRequest
		response = "destination missing"
		return status, response, values
	}

	values.ani = params.Get("ani")
	values.longMessageMode = params.Get("longMessageMode")

	return status, response, values
}

func SubmitGET(r *Route) {

	status, response, values := checkParams(r.r.URL.Query())

	if values.dnis != "" {
		// Check for special overrides on number suffix
		number := values.dnis[len(values.dnis)-4:]
		log.Printf("Number: %s", number)
		// Check for spoofing a submission error
		switch number {
		case "1400":
			status = http.StatusBadRequest
			response = "NO ROUTES"
		}
	}

	r.api.log.Fields(logger.Fields{
		"msgid":           "xxx",
		"dnis":            values.dnis,
		"ani":             values.ani,
		"message":         values.message,
		"command":         values.command,
		"longMessageMode": values.longMessageMode,
		"status":          status,
	}).Info(response)

	if status != http.StatusOK {
		r.w.Header().Set("Content-Type", "text/html")
		r.w.WriteHeader(status)
		fmt.Fprintf(r.w, response)
		return
	}

	type payload struct {
		MessageID string `json:"message_id"`
	}
	data := payload{
		MessageID: "xxx",
	}

	r.Write(data, http.StatusOK)
}
