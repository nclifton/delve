package tualet

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/burstsms/mtmo-tp/backend/sms/biz"
)

type handsetParams struct {
	message string
	dnis    string
	ani     string
}

func checkMOParams(params url.Values) (int, string, handsetParams) {
	status := http.StatusOK
	response := "spoof MO"

	values := handsetParams{}

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

	return status, response, values
}

func HandsetGET(r *Route) {

	status, response, values := checkMOParams(r.r.URL.Query())

	split, err := biz.SplitSMSParts(values.message)
	if err != nil {
		status = http.StatusBadRequest
		response = err.Error()
	}

	log.Printf("SMS Parts: %+v", split)

	uuid := uuid.New()
	MessageID := uuid.String()

	if status != http.StatusOK {
		r.w.Header().Set("Content-Type", "text/html")
		r.w.WriteHeader(status)
		fmt.Fprint(r.w, response)
		return
	}

	sarId := `"$sarId$"`
	sarParts := `"$sarParts$"`
	sarPartNumber := `"$sarPartNumber$"`
	sarIdsplit := uuid.String()

	for part, msg := range split.Parts {
		go func(part int, msg biz.Sms) {
			rand.Seed(time.Now().UnixNano())
			delay := rand.Intn(10)
			time.Sleep(time.Duration((delay + 1)) * time.Second)
			if split.CountParts > 1 {
				// ok its a multi sms, so we need valid sar parameters
				sarParts = strconv.Itoa(split.CountParts)
				sarId = sarIdsplit
				sarPartNumber = strconv.Itoa(part + 1)
			}
			r.api.log.Fields(r.r.Context(), logger.Fields{
				"msgid":   MessageID,
				"dnis":    values.dnis,
				"ani":     values.ani,
				"message": msg.Content,
				"status":  status,
			}).Info(response)

			MessageID := uuid.String()
			data := url.Values{}
			data.Set("msgid", MessageID)
			data.Set("to", values.dnis)
			data.Set("from", values.ani)
			data.Set("message", msg.Content)
			data.Set("sarId", sarId)
			data.Set("sarPartNumber", sarPartNumber)
			data.Set("sarParts", sarParts)

			req, err := http.NewRequest("POST", r.api.opts.MOEndpoint, strings.NewReader(data.Encode()))
			if err != nil {
				r.api.log.Errorf(r.r.Context(), "HandsetGET", "Could not create DLR request: %s", err)
				return
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

			resp, err := r.api.client.Do(req)
			if err != nil {
				r.api.log.Errorf(r.r.Context(), "HandsetGET", "Could not do DLR request: %s", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, _ := ioutil.ReadAll(resp.Body)
				r.api.log.Errorf(r.r.Context(), "HandsetGET", "Not OK response from %s, with code: %d, body %s", r.api.opts.DLREndpoint, resp.StatusCode, string(body))
				return
			}
		}(part, msg)
	}

}
