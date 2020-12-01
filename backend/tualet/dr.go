package tualet

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type DLRParams struct {
	To         string
	Status     string
	ReasonCode string
	MessageID  string
	MCC        string
	MNC        string
}

func (api *TualetAPI) sendDLRRequest(params *DLRParams) {

	rand.Seed(time.Now().UnixNano())

	if api.opts.DLREndpoint != "" {

		status := "DELIVRD"
		// Check for special overrides on number suffix
		number := params.To[len(params.To)-4:]
		// Check for spoofing a submission error
		switch number {
		case "1500":
			status = "ENROUTE"
		case "1507":
			status = "SENT"
		case "1501":
			status = "DELIVRD"
		case "1502":
			status = "EXPIRED"
		case "1503":
			status = "DELETED"
		case "1504":
			status = "UNDELIV"
		case "1505":
			status = "REJECTD"
		case "1506":
			status = "UNKNOWN"
		}

		go func() {
			// Introduce some random time between dlrs, so they can come out of sequence and not before they
			// have been marked sent
			delay := rand.Intn(6)
			time.Sleep(time.Duration((delay + 1)) * time.Second)
			data := url.Values{}
			data.Set("msgid", params.MessageID)
			data.Set("state", status)
			data.Set("reasoncode", params.ReasonCode)
			data.Set("to", params.To)
			data.Set("time", time.Now().UTC().Format(time.RFC3339))
			data.Set("mcc", params.MCC)
			data.Set("mnc", params.MNC)

			req, err := http.NewRequest("POST", api.opts.DLREndpoint, strings.NewReader(data.Encode()))
			if err != nil {
				api.log.Errorf("Could not create DLR request: %s", err)
				return
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

			resp, err := api.client.Do(req)
			if err != nil {
				api.log.Errorf("Could not do DLR request: %s", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, _ := ioutil.ReadAll(resp.Body)
				api.log.Errorf("Not OK response from %s, with code: %d, body %s", api.opts.DLREndpoint, resp.StatusCode, string(body))
				return
			}
		}()
	}

}
