package rpc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/burstsms/mtmo-tp/backend/lib/redis"
	"github.com/burstsms/mtmo-tp/backend/mm7/rpc/types"
	"github.com/burstsms/mtmo-tp/backend/mm7/worker"
	mms "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
)

var validMediaRegex = regexp.MustCompile(`^image\/(png|jpeg|gif)`)

func (s *MM7) Ping(p types.NoParams, r *types.PingResponse) error {
	r.Res = "PONG"
	return nil
}

func (s *MM7) Send(p types.MM7SendParams, r *types.NoReply) error {
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

func (s *MM7) ProviderSpec(p types.MM7ProviderSpecParams, r *types.MM7ProviderSpecReply) error {
	provider, err := getProviderDetail(p.ProviderKey)
	if err != nil {
		return err
	}

	r.ProviderKey = p.ProviderKey
	r.ImageSizeMaxKB = provider.ImageSizeMaxKB

	return nil
}

func (s *MM7) UpdateStatus(p types.MM7UpdateStatusParams, r *types.NoReply) error {
	return s.svc.MMS.UpdateStatus(mms.UpdateStatusParams{
		ID:          p.ID,
		MessageID:   p.MessageID,
		Status:      p.Status,
		Description: p.Description,
	})
}

func (s *MM7) DLR(p types.MM7DLRParams, r *types.NoReply) error {
	log.Printf("RPC call DLR, params: %+v", p)
	return nil
}

func (s *MM7) Deliver(p types.MM7DeliverParams, r *types.NoReply) error {
	log.Printf("RPC call Deliver, params: %+v", p)
	return nil
}

func (s *MM7) GetCachedContent(p types.MM7GetCachedContentParams, r *types.MM7GetCachedContentReply) error {
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
	contentType := http.DetectContentType(body)

	if validMediaRegex.FindString(contentType) == "" {
		return fmt.Errorf("Invalid contentType (%s) for %s", contentType, p.ContentURL)
	}

	if err := s.db.redis.Client.Set(p.ContentURL, body, time.Hour).Err(); err != nil {
		return err
	}

	r.Content = body
	return nil
}

func (s *MM7) CheckRateLimit(p types.MM7CheckRateLimitParams, r *types.MM7CheckRateLimitReply) error {
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
			ImageSizeMaxKB:  450,
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
