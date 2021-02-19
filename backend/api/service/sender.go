package service

import (
	"net/http"
	"time"

	"github.com/burstsms/mtmo-tp/backend/lib/rest"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
)

type APIResponseSender struct {
	ID             string   `json:"id"`
	AccountID      string   `json:"account_id"`
	Address        string   `json:"address"`
	MMSProviderKey string   `json:"mms_provider_key"`
	Channels       []string `json:"channels"`
	Country        string   `json:"country"`
	Comment        string   `json:"comment"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

type APIResponseSenders struct {
	Senders []APIResponseSender `json:"senders"`
}

func (s *Service) SenderListGET(hc *rest.HandlerContext) {
	account := accountFromCtx(hc)

	reply, err := s.SenderClient.FindSendersByAccountId(hc.Context(), &senderpb.FindSendersByAccountIdParams{
		AccountId: account.GetId(),
	})
	if err != nil {
		hc.LogFatal(err)
	}

	res := APIResponseSenders{
		Senders: []APIResponseSender{},
	}
	for _, s := range reply.Senders {
		res.Senders = append(res.Senders, APIResponseSender{
			ID:             s.Id,
			AccountID:      s.AccountId,
			Address:        s.Address,
			MMSProviderKey: s.MMSProviderKey,
			Channels:       s.Channels,
			Country:        s.Country,
			Comment:        s.Comment,
			CreatedAt:      s.CreatedAt.AsTime().Format(time.RFC3339),
			UpdatedAt:      s.UpdatedAt.AsTime().Format(time.RFC3339),
		})
	}

	hc.WriteJSON(res, http.StatusOK)
}
