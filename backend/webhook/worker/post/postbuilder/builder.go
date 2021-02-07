package postbuilder

import (
	"context"
	"net/http"
	"time"

	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/redis"
	"github.com/burstsms/mtmo-tp/backend/lib/workerbuilder"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/post/handler"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ClientTimeout int    `envconfig:"CLIENT_TIMEOUT"`
	RedisURL      string `envconfig:"REDIS_URL"`
}

type postService struct {
	conf    Config
	client  handler.HTTPClient
	limiter handler.Limiter
}

func NewBuilderFromEnv() *postService {
	stLog := logger.NewLogger()

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		stLog.Fatalf(context.Background(), "envconfig.Process", "failed to read env vars: %s", err)
	}

	return New(config)
}

func New(config Config) *postService {
	return &postService{conf: config}
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

	handler := handler.New(ps.client, ps.limiter)

	// TODO move health set service ready true/false into the Worker
	deps.Health.SetServiceReady(true)
	deps.Worker.Run(deps.ConsumeOptions, handler)

	return nil
}
