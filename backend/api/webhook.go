package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TestPublishOptOutWebhookRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func TestPublishOptOutWebhookPOST(r *Route) {
	account, err := r.RequireAccountContext()
	if err != nil {
		return
	}

	var req TestPublishOptOutWebhookRequest
	err = r.DecodeRequest(&req)
	if err != nil {
		return
	}

	r.api.log.Info(r.r.Context(), "TestPublishOptOutWebhookPOST", fmt.Sprint(req))

	res, err := r.api.webhook.Insert(r.r.Context(), &webhookpb.InsertParams{
		AccountId: account.ID,
		Event:     "opt_out",
		Name:      req.Name,
		URL:       req.Url,
		RateLimit: 10,
	})
	if err != nil {
		r.api.log.Error(r.r.Context(), "r.api.webhook.Insert", err.Error())
		r.WriteError("Could not process webhook", http.StatusInternalServerError)
		return
	}

	message := webhookpb.Message{
		Type:       "sms",
		Id:         "1",
		Recipient:  "recipient",
		Sender:     "name",
		Message:    req.Name,
		MessageRef: "ref",
	}

	_, err = r.api.webhook.PublishOptOut(r.r.Context(), &webhookpb.PublishOptOutParams{
		Source:        "sms_inbound",
		Timestamp:     timestamppb.New(time.Now().UTC()),
		AccountId:     account.ID,
		SourceMessage: &message,
	})
	if err != nil {
		r.api.log.Error(r.r.Context(), "r.api.webhook.PublishOptOut", err.Error())
		r.WriteError("Could not publish opt out", http.StatusInternalServerError)
		return
	}

	_, err = r.api.webhook.Delete(r.r.Context(), &webhookpb.DeleteParams{
		Id:        res.Webhook.Id,
		AccountId: account.ID,
	})
	if err != nil {
		r.api.log.Error(r.r.Context(), "r.api.webhook.Delete", err.Error())
		r.WriteError("Could not delete", http.StatusInternalServerError)
		return
	}

	r.Write(res.Webhook, http.StatusOK)
}

type WebhookCreatePOSTRequest struct {
	Event     string `json:"event" valid:"contains(link_hit|opt_out|sms_status|mms_status|sms_inbound|mms_inbound)"`
	Name      string `json:"name" valid:"length(2|100)"`
	URL       string `json:"url" valid:"webhook_url"`
	RateLimit int    `json:"rate_limit" valid:"range(0|10000)"`
}

func WebhookCreatePOST(r *Route) {
	account, err := r.RequireAccountContext()
	if err != nil {
		return
	}

	var req WebhookCreatePOSTRequest
	err = r.DecodeRequest(&req)
	if err != nil {
		return
	}

	r.api.log.Info(r.r.Context(), "WebhookCreatePOST", fmt.Sprint(req))

	res, err := r.api.webhook.Insert(r.r.Context(), &webhookpb.InsertParams{
		AccountId: account.ID,
		Event:     req.Event,
		Name:      req.Name,
		URL:       req.URL,
		RateLimit: int32(req.RateLimit),
	})
	if err != nil {
		r.api.log.Error(r.r.Context(), "r.api.webhook.Insert", err.Error())
		r.WriteError("Could not process webhook", http.StatusInternalServerError)
		return
	}

	r.Write(res.Webhook, http.StatusOK)
}

func WebhookGET(r *Route) {
	account, err := r.RequireAccountContext()
	if err != nil {
		return
	}

	webhookID := r.params.ByName("id")
	if webhookID == "" {
		r.WriteError("Invalid Param: id", http.StatusBadRequest)
		return
	}

	r.api.log.Info(r.r.Context(), "WebhookGET", webhookID)

	res, err := r.api.webhook.FindByID(r.r.Context(), &webhookpb.FindByIDParams{AccountId: account.ID, WebhookId: webhookID})
	if err != nil {
		r.api.log.Error(r.r.Context(), "r.api.webhook.get", err.Error())
		st, ok := status.FromError(err)
		if ok {
			if st.Code() == codes.NotFound {
				r.WriteError("Could not find webhook", http.StatusNotFound)
				return
			}
		}
		r.WriteError("Could not process webhook", http.StatusInternalServerError)
		return
	}

	r.Write(res.Webhook, http.StatusOK)
}

func WebhookListGET(r *Route) {
	account, err := r.RequireAccountContext()
	if err != nil {
		return
	}

	r.api.log.Info(r.r.Context(), "WebhookListGET", "*")

	res, err := r.api.webhook.Find(r.r.Context(), &webhookpb.FindParams{AccountId: account.ID})
	if err != nil {
		r.api.log.Error(r.r.Context(), "r.api.webhook.list", err.Error())
		r.WriteError("Could not find webhooks", http.StatusNotFound)
		return
	}

	r.Write(res.Webhooks, http.StatusOK)
}

func WebhookDELETE(r *Route) {
	account, err := r.RequireAccountContext()
	if err != nil {
		return
	}

	webhookID, err := strconv.Atoi(r.params.ByName("id"))
	if err != nil {
		r.WriteError("Invalid Param: id", http.StatusBadRequest)
		r.api.log.Error(r.r.Context(), "r.api.webhook.delete", err.Error())
		return
	}

	r.api.log.Info(r.r.Context(), "WebhookDELETE", fmt.Sprint(webhookID))

	_, err = r.api.webhook.Delete(r.r.Context(), &webhookpb.DeleteParams{AccountId: account.ID, Id: int64(webhookID)})
	if err != nil {
		r.api.log.Error(r.r.Context(), "r.api.webhook.delete", err.Error())
		r.WriteError("Could not delete webhook", http.StatusInternalServerError)
		return
	}

	r.WriteOK()
}

type WebhookUpdatePUTRequest struct {
	Event     string `json:"event" valid:"contains(link_hit|opt_out|sms_status|mms_status|sms_inbound|mms_inbound)"`
	Name      string `json:"name" valid:"length(2|100)"`
	URL       string `json:"url" valid:"webhook_url"`
	RateLimit int    `json:"rate_limit" valid:"range(0|10000)"`
}

func WebhookUpdatePUT(r *Route) {
	account, err := r.RequireAccountContext()
	if err != nil {
		return
	}

	var req WebhookUpdatePUTRequest
	err = r.DecodeRequest(&req)
	if err != nil {
		return
	}

	webhookID, err := strconv.Atoi(r.params.ByName("id"))
	if err != nil {
		r.WriteError("Invalid Param: id", http.StatusBadRequest)
		r.api.log.Error(r.r.Context(), "r.api.webhook.update", err.Error())
		return
	}

	r.api.log.Info(r.r.Context(), "WebhookUpdatePUT", fmt.Sprint(req))

	res, err := r.api.webhook.Update(r.r.Context(), &webhookpb.UpdateParams{
		Id:        int64(webhookID),
		AccountId: account.ID,
		Event:     req.Event,
		Name:      req.Name,
		URL:       req.URL,
		RateLimit: int32(req.RateLimit),
	})
	if err != nil {
		r.api.log.Error(r.r.Context(), "r.api.webhook.Update", err.Error())
		r.WriteError("Could not update webhook", http.StatusInternalServerError)
		return
	}

	r.Write(res.Webhook, http.StatusOK)
}
