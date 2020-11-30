package fakemm7dldrworker

import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/burstsms/mtmo-tp/backend/lib/mm7utils"
	belogger "github.com/burstsms/mtmo-tp/backend/logger"
	mm7RPC "github.com/burstsms/mtmo-tp/backend/mm7/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/mm7/worker"
)

const (
	deliver        = "deliver"
	deliveryReport = "deliveryReport"
)

var (
	deliverRegex        = regexp.MustCompile(`<DeliverReq.*>`)
	deliveryReportRegex = regexp.MustCompile(`<DeliveryReportReq.*>`)

	stripRegex = regexp.MustCompile(`>\s+<`)

	// deliverReport
	drMessageIDRegex   = regexp.MustCompile(`<MessageID.*>(.*)<\/MessageID>`)
	drStatusRegex      = regexp.MustCompile(`<MMStatus.*>(.*)<\/MMStatus>`)
	drDescriptionRegex = regexp.MustCompile(`<StatusText.*>(.*)<\/StatusText>`)

	// deliver
	dlSubjectRegex       = regexp.MustCompile(`<Subject>(.*)<\/Subject>`)
	dlRecipientRegex     = regexp.MustCompile(`<Recipients><To><Number displayOnly="false">(.*)<\/Number><\/To></Recipients>`)
	dlSenderRegex        = regexp.MustCompile(`<Sender><Number>(.*)<\/Number><\/Sender>`)
	dlTransactionIDRegex = regexp.MustCompile(`<LinkedID>(.*)<\/LinkedID>`)
)

type MM7RPCClient interface {
	Store(params mm7RPC.MM7MediaStoreParams) (r *mm7RPC.MM7MediaStoreReply, err error)
	DLR(params mm7RPC.MM7DLRParams) (r *mm7RPC.NoReply, err error)
	Deliver(params mm7RPC.MM7DeliverParams) (r *mm7RPC.NoReply, err error)
}

type FakeMM7DLDRHandler struct {
	mm7RPC MM7RPCClient
	log    *belogger.StandardLogger
}

func NewHandler(c MM7RPCClient) *FakeMM7DLDRHandler {
	return &FakeMM7DLDRHandler{
		mm7RPC: c,
		log:    belogger.NewLogger(),
	}
}

func (h *FakeMM7DLDRHandler) OnFinalFailure(body []byte) error {
	return nil
}

func (h *FakeMM7DLDRHandler) Handle(body []byte, headers map[string]interface{}) error {
	contentType, ok := headers["Content-Type"].(string)
	if !ok {
		err := errors.New("Content-Type not provided")
		h.log.Errorln(err)
		return err
	}

	if strings.HasPrefix(contentType, "application/xml") {
		bodyReq := string(body)
		bodyReq = strings.Replace(string(bodyReq), "\n", "", -1)
		bodyReq = stripRegex.ReplaceAllString(bodyReq, "><")

		action, err := getRequestAction(bodyReq)
		if err != nil {
			h.log.Errorln(err)
			return err
		}

		if action != deliveryReport {
			err = fmt.Errorf("Wrong action %s for Content/type:application/xml", action)
			h.log.Errorln(err)
			return err
		}

		if err := h.processDeliveryReport(bodyReq); err != nil {
			h.log.Errorln(err)
			return err
		}

		return nil
	}

	parts, err := mm7utils.ProcessMultiPart(contentType, bytes.NewReader(body))
	if err != nil {
		h.log.Errorln(err)
		return err
	}

	if parts == nil {
		err = errors.New("No multi parts processed")
		h.log.Errorln(err)
		return err
	}

	bodyReq := string(parts[0].Body)
	bodyReq = strings.Replace(string(bodyReq), "\n", "", -1)
	bodyReq = stripRegex.ReplaceAllString(bodyReq, "><")

	action, err := getRequestAction(bodyReq)
	if err != nil {
		h.log.Errorln(err)
		return err
	}

	if action != deliver {
		err = fmt.Errorf("Wrong action %s for Content/type:%s", action, contentType)
		h.log.Errorln(err)
		return err
	}

	if err := h.processDeliver(bodyReq, parts); err != nil {
		h.log.Errorln(err)
		return err
	}

	return nil
}

func getRequestAction(body string) (string, error) {
	if deliverRegex.MatchString(body) {
		return deliver, nil
	}

	if deliveryReportRegex.MatchString(body) {
		return deliveryReport, nil
	}

	return "", errors.New("Unknown request action")
}

func (h *FakeMM7DLDRHandler) processDeliveryReport(body string) error {
	_, err := h.mm7RPC.DLR(mm7RPC.MM7DLRParams{
		ID:          mm7utils.ExtractEntity(*drMessageIDRegex, body),
		Status:      mm7utils.ExtractEntity(*drStatusRegex, body),
		Description: mm7utils.ExtractEntity(*drDescriptionRegex, body),
	})

	return err
}

func (h *FakeMM7DLDRHandler) processDeliver(body string, parts []*mm7utils.MMSPart) error {
	params := mm7RPC.MM7DeliverParams{
		Subject:     mm7utils.ExtractEntity(*dlSubjectRegex, body),
		Sender:      mm7utils.ExtractEntity(*dlSenderRegex, body),
		Recipient:   mm7utils.ExtractEntity(*dlRecipientRegex, body),
		ProviderKey: worker.FakeProviderKey,
	}

	for _, part := range parts[1:] {
		if strings.HasPrefix(part.ContentType, "image/") {
			ext := path.Ext(part.ContentID)
			fileName := fmt.Sprintf("%s_%s", mm7utils.ExtractEntity(*dlTransactionIDRegex, body), strings.TrimSuffix(part.ContentID, ext))

			reply, err := h.mm7RPC.Store(mm7RPC.MM7MediaStoreParams{
				FileName:    fileName,
				ProviderKey: worker.FakeProviderKey,
				Data:        part.Body,
				Extension:   ext,
			})
			if err != nil {
				return err
			}

			params.ContentURLs = append(params.ContentURLs, reply.URL)
			continue
		}

		if strings.HasPrefix(part.ContentType, "text/plain") {
			params.Message = string(part.Body)
		}
	}

	_, err := h.mm7RPC.Deliver(params)
	if err != nil {
		return err
	}

	return nil
}
