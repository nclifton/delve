package rpc

import (
	"time"
)

type WebhookRecord struct {
	ID        int64     `json:"id"`
	AccountID string    `json:"account_id"`
	Event     string    `json:"event"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	RateLimit int       `json:"rate_limit"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FindParams struct {
	AccountID string
}

type FindReply struct {
	Webhooks []WebhookRecord `json:"webhooks"`
}

func (s *Webhook) Find(p FindParams, r *FindReply) error {
	var err error

	results, err := s.db.Find(p.AccountID)
	if err != nil {
		return err
	}

	webhooks := []WebhookRecord{}
	for _, w := range results {
		webhooks = append(webhooks, WebhookRecord{
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

type InsertParams struct {
	AccountID string `json:"account_id"`
	Event     string `json:"event"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	RateLimit int    `json:"rate_limit"`
}

type InsertReply struct {
	Webhook WebhookRecord
}

func (s *Webhook) Insert(p InsertParams, r *InsertReply) error {
	w, err := s.db.Insert(p.AccountID, p.Event, p.Name, p.URL, p.RateLimit)
	r.Webhook = WebhookRecord{
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

type DeleteParams struct {
	AccountID string `json:"account_id"`
	ID        string `json:"ids"`
}

func (s *Webhook) Delete(p DeleteParams, r *NoReply) error {
	err := s.db.Delete(p.AccountID, p.ID)

	return err
}

type UpdateParams struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	Event     string `json:"event"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	RateLimit int    `json:"rate_limit"`
}

type UpdateReply struct {
	Webhook WebhookRecord
}

func (s *Webhook) Update(p UpdateParams, r *UpdateReply) error {
	w, err := s.db.Update(p.ID, p.AccountID, p.Event, p.Name, p.URL, p.RateLimit)
	r.Webhook = WebhookRecord{
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

/*
type PublishContactUpdateParams struct {
	AccountID  string
	ContactID  string
	Timestamp  time.Time
	Operation  string
	OldContact struct {
		ContactRef string
		Mobile     string
		FirstName  string
		LastName   string
		Email      string
		Country    string
		Status     string
		Timezone   string
		Custom     map[string]interface{}
	}
	NewContact struct {
		ContactRef string
		Mobile     string
		FirstName  string
		LastName   string
		Email      string
		Country    string
		Status     string
		Timezone   string
		Custom     map[string]interface{}
	}
}

func (s *Webhook) PublishContactUpdate(p PublishContactUpdateParams, r *NoReply) error {
	webhooks, err := s.db.FindByEvent(p.AccountID, erpc.ContactUpdate)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return err
	}

	for _, w := range webhooks {
		err = s.db.Publish(worker.WebhookMessage{
			URL:       w.URL,
			RateLimit: w.RateLimit,
			Payload:   payload,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
*/
