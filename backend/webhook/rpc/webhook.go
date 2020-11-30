package rpc

import (
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/types"
)

func (s *Webhook) Find(p types.FindParams, r *types.FindReply) error {
	var err error

	results, err := s.db.Find(p.AccountID)
	if err != nil {
		return err
	}

	webhooks := []types.WebhookRecord{}
	for _, w := range results {
		webhooks = append(webhooks, types.WebhookRecord{
			ID:        w.ID,
			AccountID: w.AccountID,
			Event:     w.Event,
			Name:      w.Name,
			URL:       w.URL,
			RateLimit: w.RateLimit,
			CreatedAt: w.CreatedAt,
			UpdatedAt: w.UpdatedAt,
		})
	}
	r.Webhooks = webhooks
	return nil
}

func (s *Webhook) Insert(p types.InsertParams, r *types.InsertReply) error {
	w, err := s.db.Insert(p.AccountID, p.Event, p.Name, p.URL, p.RateLimit)
	r.Webhook = types.WebhookRecord{
		ID:        w.ID,
		AccountID: w.AccountID,
		Event:     w.Event,
		Name:      w.Name,
		URL:       w.URL,
		RateLimit: w.RateLimit,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}

	return err
}

func (s *Webhook) Delete(p types.DeleteParams, r *types.NoReply) error {
	err := s.db.Delete(p.AccountID, p.ID)

	return err
}

func (s *Webhook) Update(p types.UpdateParams, r *types.UpdateReply) error {
	w, err := s.db.Update(p.ID, p.AccountID, p.Event, p.Name, p.URL, p.RateLimit)
	r.Webhook = types.WebhookRecord{
		ID:        w.ID,
		AccountID: w.AccountID,
		Event:     w.Event,
		Name:      w.Name,
		URL:       w.URL,
		RateLimit: w.RateLimit,
		UpdatedAt: w.UpdatedAt,
		CreatedAt: w.CreatedAt,
	}

	return err
}
