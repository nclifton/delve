package service

import (
	"net/http"
	"time"

	"github.com/burstsms/mtmo-tp/backend/lib/rest"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
)

type APIResponseSender struct {
	Id               string   `json:"id"`
	Account_id       string   `json:"account_id"`
	Address          string   `json:"address"`
	MMS_provider_key string   `json:"mms_provider_key"`
	Channels         []string `json:"channels"`
	Country          string   `json:"country"`
	Comment          string   `json:"comment"`
	Created_at       string   `json:"created_at"`
	Updated_at       string   `json:"updated_at"`
}

type APIResponseSenders struct {
	Senders []APIResponseSender `json:"senders"`
}

func (s *Service) SenderListGET(hc *rest.HandlerContext) {
	account := accountFromCtx(hc)

	reply, err := s.SenderClient.FindSendersByAccountId(hc.Context(), &senderpb.FindSendersByAccountIdParams{
		AccountId: account.ID,
	})
	if err != nil {
		hc.LogFatal(err)
	}

	res := APIResponseSenders{
		Senders: []APIResponseSender{},
	}
	for _, s := range reply.Senders {
		res.Senders = append(res.Senders, APIResponseSender{
			Id:               s.Id,
			Account_id:       s.AccountId,
			Address:          s.Address,
			MMS_provider_key: s.MMSProviderKey,
			Channels:         s.Channels,
			Country:          s.Country,
			Comment:          s.Comment,
			Created_at:       s.CreatedAt.AsTime().Format(time.RFC3339),
			Updated_at:       s.UpdatedAt.AsTime().Format(time.RFC3339),
		})
	}

	hc.WriteJSON(res, http.StatusOK)
}
