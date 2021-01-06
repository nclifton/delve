package api

import (
	"log"
	"net/http"
	"time"

	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
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

	res, err := r.api.webhook.Insert(r.r.Context(), &webhookpb.InsertParams{
		AccountId: account.ID,
		Event:     "opt_out",
		Name:      req.Name,
		URL:       req.Url,
		RateLimit: 10,
	})
	if err != nil {
		log.Printf("Could not call insert webhook: %s", err.Error())
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
		log.Printf("Could not call publish opt out for webhook: %s", err.Error())
		r.WriteError("Could not publish opt out", http.StatusInternalServerError)
		return
	}

	_, err = r.api.webhook.Delete(r.r.Context(), &webhookpb.DeleteParams{
		Id:        res.Webhook.Id,
		AccountId: account.ID,
	})
	if err != nil {
		log.Printf("Could not call delete for webhook: %s", err.Error())
		r.WriteError("Could not delete", http.StatusInternalServerError)
		return
	}

	r.Write(res.Webhook, http.StatusOK)
}
