package rpc

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/burstsms/mtmo-tp/backend/lib/redis"
	"github.com/burstsms/mtmo-tp/backend/mm7/worker"
	mms "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
)

type PingResponse struct {
	Res string
}

func (s *MM7) Ping(p NoParams, r *PingResponse) error {
	r.Res = "PONG"
	return nil
}

type MM7SendParams struct {
	ID          string
	Subject     string
	Message     string
	Sender      string
	Recipient   string
	ContentURLs []string
	ProviderKey string
}

func (s *MM7) Send(p MM7SendParams, r *NoReply) error {
	msg := worker.SubmitMessage{
		ID:          p.ID,
		Subject:     p.Subject,
		Message:     p.Message,
		Sender:      p.Sender,
		Recipient:   p.Recipient,
		ContentURLs: p.ContentURLs,
	}

	provider, err := getProviderDetail(p.ProviderKey)
	if err != nil {
		return err
	}

	return s.db.Publish(msg, provider.QueueNameSubmit)
}

type MM7ProviderSpecParams struct {
	ProviderKey string
}

type MM7ProviderSpecReply struct {
	ProviderKey    string
	ImageSizeMaxKB int
}

func (s *MM7) ProviderSpec(p MM7ProviderSpecParams, r *MM7ProviderSpecReply) error {
	provider, err := getProviderDetail(p.ProviderKey)
	if err != nil {
		return err
	}

	r.ProviderKey = p.ProviderKey
	r.ImageSizeMaxKB = provider.ImageSizeMaxKB

	return nil
}

type MM7UpdateStatusParams struct {
	ID          string
	MessageID   string
	Status      string
	Description string
}

func (s *MM7) UpdateStatus(p MM7UpdateStatusParams, r *NoReply) error {
	return s.svc.MMS.UpdateStatus(mms.UpdateStatusParams{
		ID:          p.ID,
		MessageID:   p.MessageID,
		Status:      p.Status,
		Description: p.Description,
	})
}

type MM7DLRParams struct {
	ID          string
	Status      string
	Description string
}

func (s *MM7) DLR(p MM7DLRParams, r *NoReply) error {
	log.Printf("RPC call DLR, params: %+v", p)
	return nil
}

type MM7DeliverParams struct {
	Subject     string
	Message     string
	Sender      string
	Recipient   string
	ContentURLs []string
	ProviderKey string
}

func (s *MM7) Deliver(p MM7DeliverParams, r *NoReply) error {
	log.Printf("RPC call Deliver, params: %+v", p)
	return nil
}

type MM7GetCachedContentParams struct {
	ContentURL string
}

type MM7GetCachedContentReply struct {
	Content []byte
}

func (s *MM7) GetCachedContent(p MM7GetCachedContentParams, r *MM7GetCachedContentReply) error {
	image, err := s.db.redis.Client.Get(p.ContentURL).Result()
	if err != nil && err != redis.RedisNil {
		return err
	}

	if image != "" {
		r.Content = []byte(image)
		return nil
	}

	imageRes, err := http.Get(p.ContentURL)
	if err != nil {
		return err
	}

	defer imageRes.Body.Close()

	body, err := ioutil.ReadAll(imageRes.Body)
	if err != nil {
		return err
	}

	if err := s.db.redis.Client.Set(p.ContentURL, body, time.Hour).Err(); err != nil {
		return err
	}

	r.Content = body
	return nil
}

type MM7CheckRateLimitParams struct {
	ProviderKey string
}

type MM7CheckRateLimitReply struct {
	Allow bool
}

func (s *MM7) CheckRateLimit(p MM7CheckRateLimitParams, r *MM7CheckRateLimitReply) error {
	provider, err := getProviderDetail(p.ProviderKey)
	if err != nil {
		return err
	}

	r.Allow = s.db.limiter.Allow(p.ProviderKey, provider.Rate, provider.Burst)

	return nil
}

type ProviderDetail struct {
	Rate            float64
	Burst           int
	ImageSizeMaxKB  int
	QueueNameSubmit string
}

func getProviderDetail(providerKey string) (ProviderDetail, error) {
	switch providerKey {
	case worker.FakeProviderKey:
		return ProviderDetail{
			Rate:            1.0,
			Burst:           20,
			ImageSizeMaxKB:  400,
			QueueNameSubmit: worker.QueueNameSubmitFake,
		}, nil
	case worker.OptusProviderKey:
		return ProviderDetail{
			Rate:            1.0,
			Burst:           20,
			ImageSizeMaxKB:  400,
			QueueNameSubmit: worker.QueueNameSubmitOptus,
		}, nil
	case worker.MgageProviderKey:
		return ProviderDetail{
			Rate:            1.0,
			Burst:           20,
			ImageSizeMaxKB:  400,
			QueueNameSubmit: worker.QueueNameSubmitMgage,
		}, nil
	default:
		return ProviderDetail{}, errors.New("Unknown provider key")
	}
}
