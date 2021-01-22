package postbuilder

import (
	"context"
	"net/http"
	"time"

	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/redis"
	"github.com/burstsms/mtmo-tp/backend/lib/workerbuilder"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/post/handler"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ClientTimeout int    `envconfig:"CLIENT_TIMEOUT"`
	RedisURL      string `envconfig:"REDIS_URL"`

	RabbitExchange        string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType    string `envconfig:"RABBIT_EXCHANGE_TYPE"`
	RabbitPrefetchedCount int    `envconfig:"RABBIT_PREFETCHED_COUNT"`
}

type postService struct {
	conf    Config
	client  handler.HTTPClient
	limiter handler.Limiter
}

func NewBuilderFromEnv() *postService {
	stLog := logger.NewLogger()

	var config Config
	if err := envconfig.Process("webhook", &config); err != nil {
		stLog.Fatalf(context.Background(), "envconfig.Process", "failed to read env vars: %s", err)
	}

	return New(config)
}

func New(config Config) *postService {
	return &postService{conf: config}
}

func (wb *postService) SetClient(client handler.HTTPClient) {
	wb.client = client
}

func (wb *postService) SetLimiter(limiter handler.Limiter) {
	wb.limiter = limiter
}

func (wb *postService) Run(deps workerbuilder.Deps) error {
	if wb.client == nil {
		wb.client = &http.Client{
			Timeout: time.Duration(wb.conf.ClientTimeout) * time.Second,
		}
	}

	if wb.limiter == nil {
		limiter, err := redis.NewLimiter(wb.conf.RedisURL)
		if err != nil {
			return err
		}
		wb.limiter = limiter
	}

	opts := rabbit.ConsumeOptions{
		PrefetchCount: wb.conf.RabbitPrefetchedCount,
		Exchange:      wb.conf.RabbitExchange,
		ExchangeType:  wb.conf.RabbitExchangeType,
		QueueName:     "webhook",
		RetryScale:    rabbit.RetryScale,
	}

	handler := handler.New(wb.client, wb.limiter)

	deps.Worker.Run(opts, handler)

	return nil
}
