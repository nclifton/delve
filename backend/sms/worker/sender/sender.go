package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	alaris "github.com/burstsms/mtmo-tp/backend/lib/alaris/client"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	rpc "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/sms/worker/msg"
)

type SenderHandler struct {
	smsRPC       *rpc.Client
	http         *http.Client
	limiter      Limiter
	alarisClient *alaris.Service
}

type Limiter interface {
	Allow(url string, rate float64, burst int) bool
}

type AlarisSingleResponse struct {
	MessageID  string `json:"message_id"`
	DNIS       string `json:"dnis"`
	SegmentNum string `json:"segment_num"`
}

func NewHandler(c *rpc.Client, http *http.Client, limiter Limiter, alarisClient *alaris.Service) *SenderHandler {
	return &SenderHandler{
		smsRPC:       c,
		http:         http,
		limiter:      limiter,
		alarisClient: alarisClient,
	}
}

func (h *SenderHandler) OnFinalFailure(ctx context.Context, body []byte) error {

	jobdata := &msg.SMSSendMessageSpec{}
	err := json.NewDecoder(bytes.NewReader(body)).Decode(&jobdata)
	if err != nil {
		return rabbit.NewErrWorkerMessageParse(err.Error())
	}

	// update sms status
	err = h.smsRPC.MarkFailed(rpc.MarkFailedParams{ID: jobdata.ID, AccountID: jobdata.AccountID})
	if err != nil {
		log.Printf("[SMS Send] error marking %s as failed: %s", jobdata.ID, err)
		return err
	}
	return nil
}

func (h *SenderHandler) Handle(ctx context.Context, body []byte, headers map[string]interface{}) error {

	jobdata := &msg.SMSSendMessageSpec{}
	err := json.NewDecoder(bytes.NewReader(body)).Decode(&jobdata)
	if err != nil {
		return rabbit.NewErrWorkerMessageParse(err.Error())
	}
	log.Printf("[SMS Send] got message: %+v", jobdata)
	if !h.limiter.Allow(jobdata.AlarisUser, float64(3000), 3500) {
		log.Printf("[SMS Send] retrying %s due to ratelimit", jobdata.ID)
		return rabbit.NewErrRetryWorkerMessage(fmt.Sprintf("[SMS Send] retrying %s due to ratelimit", jobdata.ID))
	}

	// hand off to alaris
	messageID, err := h.alarisClient.SendSMS(alaris.SendSMSParams{
		Username:        jobdata.AlarisUser,
		Password:        jobdata.AlarisPass,
		Command:         "submit",
		Message:         jobdata.Message,
		DNIS:            jobdata.Recipient,
		ANI:             jobdata.Sender,
		LongMessageMode: "split",
		URL:             jobdata.AlarisURL,
	})
	if err != nil {
		aerror, ok := err.(*alaris.AlarisClientError)
		if ok {
			if aerror.RetryAble {
				return rabbit.NewErrRetryWorkerMessage(fmt.Sprintf("[SMS Send] Failed sending sms to alaris Error: %s", aerror.Error()))
			}
		}
		return fmt.Errorf("[SMS Send] Failed sending sms to alaris Error: %s", aerror.Error())
	}

	log.Printf("[SMS Send] sent msg to alaris SMSID:(%s), MessageID:(%s)", jobdata.ID, messageID)

	// update sms status
	err = h.smsRPC.MarkSent(rpc.MarkSentParams{ID: jobdata.ID, MessageID: messageID, AccountID: jobdata.AccountID})
	if err != nil {
		log.Printf("[SMS Send] error marking %s as sent: %s", jobdata.ID, err)
		return err
	}

	return nil
}
