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
	env                         *WebhookEnv
	nrOpts                      *agent.Options
	queueConnection             rabbit.Conn
	httpClient                  *http.Client
	limiter                     *redis.Limiter
	handler                     *handler.Webhook
	queueWorker                 *rabbit.Worker
	rcOpts                      rabbit.ConsumeOptions
	ignoreClosedQueueConnection bool
}

func New() *worker {
	w := worker{
		env:                         &WebhookEnv{},
		ignoreClosedQueueConnection: false,
	}
	w.loadEnv()

	return &w
}

func (w *worker) Run() {
	w.prepare()
	w.queueWorker.Run(w.rcOpts, w.handler)
}

func (w *worker) IgnoreClosedQueueConnection() {
	w.ignoreClosedQueueConnection = true
}

func (w *worker) SetEnv(env *WebhookEnv) {
	w.env = env
}

func (w *worker) SetNrOpts(nrOpts *agent.Options) {
	w.nrOpts = nrOpts
}

func (w *worker) loadEnv() {
	if err := envconfig.Process("", w.env); err != nil {
		log.Fatal("failed to read env vars:", err)
	}
}

func (w *worker) prepare() {
	s, _ := json.MarshalIndent(w.env, "", "\t")
	log.Printf("ENV: %s", s)

	w.prepareDependencies()

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
}

func (w *worker) prepareDependencies() {
	// new relic options
	if w.nrOpts == nil {
		w.SetNrOpts(&agent.Options{
			AppName:                  w.env.NRName,
			NewRelicLicense:          w.env.NRLicense,
			DistributedTracerEnabled: w.env.NRTracing,
		})
	}
	if w.queueConnection == nil {
		w.createQueueConnection()
	}
	if w.httpClient == nil {
		w.createHttpClient()
	}
	if w.limiter == nil {
		w.createLimiter()
	}
}

func (w *worker) createQueueConnection() {
	var err error
	w.queueConnection, err = rabbit.Connect(w.env.RabbitURL, w.ignoreClosedQueueConnection)
	if err != nil {
		log.Fatalf("Failed to initialise queue worker: %s reason: %s\n", w.env.WorkerQueueName, err)
	}
}

func (w *worker) createHttpClient() {
	w.httpClient = &http.Client{
		Timeout: time.Duration(w.env.ClientTimeout) * time.Second,
	}
}

func (w *worker) createLimiter() {
	var err error
	w.limiter, err = redis.NewLimiter(w.env.RedisURL)
	if err != nil {
		log.Fatalf("Failed to initialise queue worker: %s reason: %s\n", w.env.WorkerQueueName, err)
	}
}
