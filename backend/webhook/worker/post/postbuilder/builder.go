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

	//TODO revise worker and post builder to move more configuration and dependency setup to the lib/workerbuilder package

	WorkerName string `envconfig:"WORKER_NAME" default:"webhook-post-worker"`                                   // no longer required - use ContainerName
	RabbitURL  string `envconfig:"RABBIT_URL" default:"amqp://tp:TheToiletPaperPassword@rabbitmq:5672/webhook"` // no longer required here

	ClientTimeout int    `envconfig:"CLIENT_TIMEOUT"`
	RedisURL      string `envconfig:"REDIS_URL"`

	// TODO, these are generic environment variables for a worker in it's own environment, should move to the worker builder deps
	RabbitExchange        string `envconfig:"RABBIT_EXCHANGE"`         // no longer required here
	RabbitExchangeType    string `envconfig:"RABBIT_EXCHANGE_TYPE"`    // no longer required here
	RabbitPrefetchedCount int    `envconfig:"RABBIT_PREFETCHED_COUNT"` // no longer required here
}

type postService struct {
	conf    Config
	client  handler.HTTPClient
	limiter handler.Limiter
}

const name string = "webhook"

func NewBuilderFromEnv() *postService {
	stLog := logger.NewLogger()

	var config Config
	if err := envconfig.Process(name, &config); err != nil {
		stLog.Fatalf(context.Background(), "envconfig.Process", "failed to read env vars: %s", err)
	}

	return New(config)
}

func New(config Config) *postService {
	return &postService{conf: config}
}

func (ps *postService) WorkerName() string {
	return ps.conf.WorkerName
}

func (ps *postService) RabbitURL() string {
	return ps.conf.RabbitURL
}

func (ps *postService) SetClient(client handler.HTTPClient) {
	ps.client = client
}

func (ps *postService) SetLimiter(limiter handler.Limiter) {
	ps.limiter = limiter
}

func (ps *postService) Run(deps workerbuilder.Deps) error {

	if ps.client == nil {
		ps.client = &http.Client{
			Timeout: time.Duration(ps.conf.ClientTimeout) * time.Second,
		}
	}

	if ps.limiter == nil {
		limiter, err := redis.NewLimiter(ps.conf.RedisURL)
		if err != nil {
			return err
		}
		ps.limiter = limiter
	}

	opts := rabbit.ConsumeOptions{
		PrefetchCount:        ps.conf.RabbitPrefetchedCount,
		Exchange:             ps.conf.RabbitExchange,
		ExchangeType:         ps.conf.RabbitExchangeType,
		QueueName:            name,
		RetryScale:           rabbit.RetryScale,
		AllowConnectionClose: deps.AllowConnectionClose,
	}

	handler := handler.New(ps.client, ps.limiter)

	deps.Health.SetServiceReady(true)
	deps.Worker.Run(opts, handler)

	return nil
}
