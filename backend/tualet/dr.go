package tualet

import (
	"io/ioutil"
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

	if api.opts.DLREndpoint != "" {
		data := url.Values{}
		data.Set("msgid", params.MessageID)
		data.Set("state", params.Status)
		data.Set("reasoncode", params.ReasonCode)
		data.Set("to", params.To)
		data.Set("time", time.Now().UTC().Format(time.RFC3339))
		data.Set("mcc", params.MCC)
		data.Set("mnc", params.MNC)

		req, err := http.NewRequest("POST", api.opts.DLREndpoint, strings.NewReader(data.Encode()))
		if err != nil {
			api.log.Errorf("Could not create DLR request: %s", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

		resp, err := api.client.Do(req)
		if err != nil {
			api.log.Errorf("Could not do DLR request: %s", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			api.log.Errorf("Not OK response from %s, with code: %d, body %s", api.opts.DLREndpoint, resp.StatusCode, string(body))
		}

	}

}
