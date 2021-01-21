package api

import (
	"net/http"
	"time"

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

func SenderListGET(r *Route) {
	account, err := r.RequireAccountContext()
	if err != nil {
		return
	}

	r.api.log.Info(r.r.Context(), "SenderListGET", "*")

	rpcReply, err := r.api.sender.FindByAccountId(r.r.Context(), &senderpb.FindByAccountIdParams{
		AccountId: account.ID,
	})
	if err != nil {
		r.api.log.Error(r.r.Context(), "r.api.sender.FindByAccountId", err.Error())
		r.WriteError("Could not find senders", http.StatusNotFound)
		return
	}

	res := APIResponseSenders{
		Senders: []APIResponseSender{},
	}
	for _, s := range rpcReply.Senders {
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

	r.Write(res, http.StatusOK)
}
