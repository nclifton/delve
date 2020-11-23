package main

import (
	"log"
	"net/http"
	"time"

	agent "github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/redis"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/msg"
	handler "github.com/burstsms/mtmo-tp/backend/webhook/worker/post"
	"github.com/kelseyhightower/envconfig"
)

type WebhookEnv struct {
	RPCPort         int    `envconfig:"RPC_PORT"`
	RPCHost         string `envconfig:"RPC_HOST"`
	RabbitURL       string `envconfig:"RABBIT_URL"`
	ClientTimeout   int    `envconfig:"CLIENT_TIMEOUT"`
	WorkerQueueName string `envconfig:"WORKER_QUEUE_NAME"`
	RedisURL        string `envconfig:"REDIS_URL"`
}

func main() {
	var env WebhookEnv
	err := envconfig.Process("webhook", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	rabbitmq, err := rabbit.Connect(env.RabbitURL)
	if err != nil {
		log.Fatalf("failed to initialise rabbit worker: %s reason: %s\n", env.WorkerQueueName, err)
	}

	client := &http.Client{
		Timeout: time.Duration(env.ClientTimeout) * time.Second,
	}

	limiter, err := redis.NewLimiter(env.RedisURL)
	if err != nil {
		log.Fatalf("failed to initialise rabbit worker: %s reason: %s\n", env.WorkerQueueName, err)
	}

	wHandler := handler.NewHandler(client, limiter)

	opts := rabbit.ConsumeOptions{
		PrefetchCount: 1,
		QueueName:     msg.WebhookMessage.Queue,
		Exchange:      msg.WebhookMessage.Exchange,
		ExchangeType:  msg.WebhookMessage.ExchangeType,
		RouteKey:      msg.WebhookMessage.RouteKey,
		RetryScale:    rabbit.RetryScale,
	}

	worker := rabbit.NewWorker(opts.QueueName, rabbitmq, &agent.Options{})
	worker.Run(opts, wHandler)
}
