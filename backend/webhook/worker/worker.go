package worker

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"

	agent "github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/redis"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/msg"
	handler "github.com/burstsms/mtmo-tp/backend/webhook/worker/post"
)

type WebhookEnv struct {
	RPCPort         int    `envconfig:"RPC_PORT"`
	RPCHost         string `envconfig:"RPC_HOST"`
	RabbitURL       string `envconfig:"RABBIT_URL"`
	ClientTimeout   int    `envconfig:"CLIENT_TIMEOUT"`
	WorkerQueueName string `envconfig:"WORKER_QUEUE_NAME"`
	RedisURL        string `envconfig:"REDIS_URL"`

	NRName    string `envconfig:"NR_NAME"`
	NRLicense string `envconfig:"NR_LICENSE"`
	NRTracing bool   `envconfig:"NR_TRACING"`
}

type worker struct {
	Env                         *WebhookEnv
	nrOpts                      *agent.Options
	queueConnection             rabbit.Conn
	httpClient                  *http.Client
	limiter                     *redis.Limiter
	handler                     *handler.Webhook
	queueWorker                 *rabbit.Worker
	rcOpts                      rabbit.ConsumeOptions
	IgnoreClosedQueueConnection bool
}

func New() *worker {
	w := worker{
		Env:                         &WebhookEnv{},
		IgnoreClosedQueueConnection: false,
	}
	if err := envconfig.Process("", w.Env); err != nil {
		log.Fatal("failed to read env vars:", err)
	}
	return &w
}

func (w *worker) Run() {
	s, _ := json.MarshalIndent(w.Env, "", "\t")
	log.Printf("ENV: %s", s)

	// new relic options
	if w.nrOpts == nil {
		w.nrOpts = &agent.Options{
			AppName:                  w.Env.NRName,
			NewRelicLicense:          w.Env.NRLicense,
			DistributedTracerEnabled: w.Env.NRTracing,
		}
	}
	var err error
	if w.queueConnection == nil {
		w.queueConnection, err = rabbit.Connect(w.Env.RabbitURL, w.IgnoreClosedQueueConnection)
		if err != nil {
			log.Fatalf("Failed to initialise queue worker: %s reason: %s\n", w.Env.WorkerQueueName, err)
		}
	}
	if w.httpClient == nil {
		w.httpClient = &http.Client{
			Timeout: time.Duration(w.Env.ClientTimeout) * time.Second,
		}
	}
	if w.limiter == nil {
		w.limiter, err = redis.NewLimiter(w.Env.RedisURL)
		if err != nil {
			log.Fatalf("Failed to initialise queue worker: %s reason: %s\n", w.Env.WorkerQueueName, err)
		}
	}
	w.handler = handler.NewHandler(w.httpClient, w.limiter)
	w.rcOpts = rabbit.ConsumeOptions{
		PrefetchCount: 1,
		QueueName:     msg.WebhookMessage.Queue,
		Exchange:      msg.WebhookMessage.Exchange,
		ExchangeType:  msg.WebhookMessage.ExchangeType,
		RouteKey:      msg.WebhookMessage.RouteKey,
		RetryScale:    rabbit.RetryScale,
	}
	w.queueWorker = rabbit.NewWorker(w.rcOpts.QueueName, w.queueConnection, w.nrOpts)
	w.queueWorker.Run(w.rcOpts, w.handler)
}
