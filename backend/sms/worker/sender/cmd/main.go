package main

import (
	"log"
	"net/http"
	"time"

	alaris "github.com/burstsms/mtmo-tp/backend/lib/alaris/client"
	agent "github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/redis"
	rpc "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/sms/worker/msg"
	sender "github.com/burstsms/mtmo-tp/backend/sms/worker/sender"
	"github.com/kelseyhightower/envconfig"
)

// TODO put this in config
var workerName = "sender"

type Env struct {
	RabbitURL     string `envconfig:"RABBIT_URL"`
	SMSRPCAddress string `envconfig:"SMS_RPC_ADDRESS"`
	ClientTimeout int    `envconfig:"CLIENT_TIMEOUT"`
	RedisURL      string `envconfig:"REDIS_URL"`
	AlarisURL     string `envconfig:"ALARIS_URL"`

	NRName    string `envconfig:"NR_NAME"`
	NRLicense string `envconfig:"NR_LICENSE"`
	NRTracing bool   `envconfig:"NR_TRACING"`
}

func main() {
	log.Println("Starting service...")

	var env Env
	err := envconfig.Process("sms", &env)
	if err != nil {
		log.Fatal("Failed to read env vars:", err)
	}

	log.Printf("ENV: %+v", env)

	// TODO service/worker level config for this url
	rabbitmq, err := rabbit.Connect(env.RabbitURL)
	if err != nil {
		log.Fatalf("Failed to initialise rabbit worker: %s reason: %s\n", workerName, err)
	}

	// TODO put this data in config
	opts := rabbit.ConsumeOptions{
		PrefetchCount: 1,
		QueueName:     msg.SMSSendMessage.Queue,
		Exchange:      msg.SMSSendMessage.Exchange,
		ExchangeType:  msg.SMSSendMessage.ExchangeType,
		RouteKey:      msg.SMSSendMessage.RouteKey,
		RetryScale:    []time.Duration{5 * time.Second, 10 * time.Second, 20 * time.Second},
	}

	nrOpts := &agent.Options{
		AppName:                  env.NRName,
		NewRelicLicense:          env.NRLicense,
		DistributedTracerEnabled: env.NRTracing,
	}

	worker := rabbit.NewWorker(workerName, rabbitmq, nrOpts)
	http := &http.Client{
		Timeout: time.Duration(env.ClientTimeout) * time.Second,
	}

	limiter, err := redis.NewLimiter(env.RedisURL)
	if err != nil {
		log.Fatalf("Failed to initialise rabbit worker: %s reason: %s\n", workerName, err)
	}

	client := rpc.New(env.SMSRPCAddress)

	alarisClient, err := alaris.NewService(env.AlarisURL, http)
	if err != nil {
		log.Fatalf("Failed to initialise rabbit worker: %s reason: %s\n", workerName, err)
	}

	log.Println("Service started")
	worker.Run(opts, sender.NewHandler(client, http, limiter, alarisClient))
}
