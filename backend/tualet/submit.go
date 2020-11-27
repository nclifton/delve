package tualet

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/burstsms/mtmo-tp/backend/logger"
	"github.com/burstsms/mtmo-tp/backend/sms/biz"
	"github.com/google/uuid"
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

	dlrStatus := `DELIVRD`
	dlrCode := "000"

	if values.dnis != "" {
		// Check for special overrides on number suffix
		number := values.dnis[len(values.dnis)-4:]
		// Check for spoofing a submission error
		switch number {
		case "1400":
			status = http.StatusBadRequest
			response = "NO ROUTES"
		}
	}

	count, err := biz.IsValidSMS(values.message, biz.SMSOptions{MaxParts: 8})
	if err != nil {
		status = http.StatusBadRequest
		response = err.Error()
	}

	uuid := uuid.New()
	MessageID := uuid.String()

	r.api.log.Fields(logger.Fields{
		"msgid":           MessageID,
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
		fmt.Fprint(r.w, response)
		return
	}

	if count > 1 {
		// ok its a multi sms, so we need a multi response
		type payload struct {
			MessageID  string `json:"message_id"`
			DNIS       string `json:"dnis"`
			SegmentNum string `json:"segment_num"`
		}
		data := []payload{}

		for segment := 1; segment <= count; segment++ {
			MessageID := uuid.String()
			data = append(data, payload{
				MessageID:  MessageID,
				DNIS:       values.dnis,
				SegmentNum: strconv.Itoa(segment),
			})
			dlrParams := DLRParams{
				To:         values.dnis,
				Status:     dlrStatus,
				ReasonCode: dlrCode,
				MessageID:  MessageID,
				MCC:        `61`,
				MNC:        `6142`,
			}
			r.api.sendDLRRequest(&dlrParams)

		}

		r.Write(data, http.StatusOK)

	} else {

		dlrParams := DLRParams{
			To:         values.dnis,
			Status:     dlrStatus,
			ReasonCode: dlrCode,
			MessageID:  MessageID,
			MCC:        `61`,
			MNC:        `6142`,
		}

		type payload struct {
			MessageID string `json:"message_id"`
		}
		data := payload{
			MessageID: MessageID,
		}

		r.Write(data, http.StatusOK)

		r.api.sendDLRRequest(&dlrParams)
	}

}
