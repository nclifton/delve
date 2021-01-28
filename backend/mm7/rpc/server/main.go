package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/lib/s3"

	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/redis"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	mm7RPC "github.com/burstsms/mtmo-tp/backend/mm7/rpc"
	mms "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	log.Println("Starting service...")

	var env mm7RPC.Env
	err := envconfig.Process("mm7", &env)
	if err != nil {
		log.Fatal("Failed to read env vars:", err)
	}

	log.Printf("ENV: %+v", env)

	// Register service with New Relic
	nr.CreateApp(&nr.Options{
		AppName:                  env.NRName,
		NewRelicLicense:          env.NRLicense,
		DistributedTracerEnabled: env.NRTracing,
	})

	rabbitmq, err := rabbit.Connect(env.RabbitURL)
	if err != nil {
		log.Fatalf("Failed to initialise rabbitmq: %s reason: %s\n", mm7RPC.Name, err)
	}

	redisCon, err := redis.Connect(env.RedisURL)
	if err != nil {
		log.Fatalf("Failed to initialise redis: %s reason: %s\n", mm7RPC.Name, err)
	}

	limiter, err := redis.NewLimiter(env.RedisURL)
	if err != nil {
		log.Fatalf("Failed to initialise redis: %s reason: %s\n", mm7RPC.Name, err)
	}

	s3ServiceParams := s3.AWSServiceS3Params{
		Region:    env.AWSRegion,
		AccessKey: env.AWSAccessKey,
		SecretKey: env.AWSSecretKey,
		PublicUrl: env.AWSS3PublicUrl,
		PathStyle: env.AWSS3PathStyle,
	}

	s3Svc, err := s3.NewService(s3ServiceParams)
	if err != nil {
		log.Fatalf("Failed to initialise s3: %s reason: %s\n", mm7RPC.Name, err)
	}

	port := env.ContainerPort

	rabbitOpts := mm7RPC.RabbitPublishOptions{
		Exchange:     env.RabbitExchange,
		ExchangeType: env.RabbitExchangeType,
	}

	configVar := mm7RPC.ConfigVar{
		AWSRegion:      env.AWSRegion,
		MMSMediaBucket: env.MMSMediaBucket,
	}

	svc := mm7RPC.ConfigSvc{
		S3:  s3Svc,
		MMS: mms.New(env.MMSHost),
	}

	server, err := rpc.NewServer(mm7RPC.NewService(rabbitmq, rabbitOpts, redisCon, limiter, svc, configVar), port)
	if err != nil {
		log.Fatalf("Failed to initialise service: %s reason: %s\n", mm7RPC.Name, err)
	}

	log.Printf("%s service initialised and available on port %d", mm7RPC.Name, port)
	server.Listen()
}
