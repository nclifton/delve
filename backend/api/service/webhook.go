package service

import (
	"net/http"
	"strconv"

	"github.com/burstsms/mtmo-tp/backend/lib/rest"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WebhookCreatePOSTRequest struct {
	Event     string `json:"event" valid:"contains(link_hit|opt_out|sms_status|mms_status|sms_inbound|mms_inbound)"`
	Name      string `json:"name" valid:"length(2|100)"`
	URL       string `json:"url" valid:"webhook_url"`
	RateLimit int    `json:"rate_limit" valid:"range(0|10000)"`
}

func (s *Service) WebhookCreatePOST(hc *rest.HandlerContext) {
	account := accountFromCtx(hc)

	var req WebhookCreatePOSTRequest
	if err := hc.DecodeJSON(&req); err != nil {
		return
	}

	res, err := s.WebhookClient.Insert(hc.Context(), &webhookpb.InsertParams{
		AccountId: account.ID,
		Event:     req.Event,
		Name:      req.Name,
		URL:       req.URL,
		RateLimit: int32(req.RateLimit),
	})
	if err != nil {
		hc.LogFatal(err)
	}

	hc.WriteJSON(res.Webhook, http.StatusOK)
}

func (s *Service) WebhookGET(hc *rest.HandlerContext) {
	account := accountFromCtx(hc)

	id := hc.Params().ByName("id")
	if id == "" {
		hc.WriteJSONError("Invalid Param: id", http.StatusBadRequest, nil)
		return
	}

	res, err := s.WebhookClient.FindByID(hc.Context(), &webhookpb.FindByIDParams{AccountId: account.ID, WebhookId: id})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			if st.Code() == codes.NotFound {
				hc.WriteJSONError("Not found", http.StatusNotFound, err)
				return
			}
		}
		hc.LogFatal(err)
	}

	hc.WriteJSON(res.Webhook, http.StatusOK)
}

func (s *Service) WebhookListGET(hc *rest.HandlerContext) {
	account := accountFromCtx(hc)

	res, err := s.WebhookClient.Find(hc.Context(), &webhookpb.FindParams{AccountId: account.ID})
	if err != nil {
		hc.LogFatal(err)
	}

	hc.WriteJSON(res.Webhooks, http.StatusOK)
}

func (s *Service) WebhookDELETE(hc *rest.HandlerContext) {
	account := accountFromCtx(hc)

	id, err := strconv.Atoi(hc.Params().ByName("id"))
	if err != nil {
		hc.WriteJSONError("Invalid Param: id", http.StatusBadRequest, nil)
		return
	}

	_, err = s.WebhookClient.Delete(hc.Context(), &webhookpb.DeleteParams{AccountId: account.ID, Id: int64(id)})
	if err != nil {
		hc.LogFatal(err)
	}

	hc.WriteJSONSuccess("Deleted")
}

type WebhookUpdatePUTRequest struct {
	Event     string `json:"event" valid:"contains(link_hit|opt_out|sms_status|mms_status|sms_inbound|mms_inbound)"`
	Name      string `json:"name" valid:"length(2|100)"`
	URL       string `json:"url" valid:"webhook_url"`
	RateLimit int    `json:"rate_limit" valid:"range(0|10000)"`
}

func (s *Service) WebhookUpdatePUT(hc *rest.HandlerContext) {
	account := accountFromCtx(hc)

	id, err := strconv.Atoi(hc.Params().ByName("id"))
	if err != nil {
		hc.WriteJSONError("Invalid Param: id", http.StatusBadRequest, nil)
		return
	}

	var req WebhookUpdatePUTRequest
	if err := hc.DecodeJSON(&req); err != nil {
		return
	}

	res, err := s.WebhookClient.Update(hc.Context(), &webhookpb.UpdateParams{
		Id:        int64(id),
		AccountId: account.ID,
		Event:     req.Event,
		Name:      req.Name,
		URL:       req.URL,
		RateLimit: int32(req.RateLimit),
	})
	if err != nil {
		hc.LogFatal(err)
	}

	hc.WriteJSON(res.Webhook, http.StatusOK)
}
