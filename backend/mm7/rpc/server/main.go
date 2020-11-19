package main

import (
	"log"

	"github.com/burstsms/mtmo-tp/backend/lib/s3"

	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/redis"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	mm7RPC "github.com/burstsms/mtmo-tp/backend/mm7/rpc"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	var env mm7RPC.Env
	err := envconfig.Process("mm7", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}

	rabbitmq, err := rabbit.Connect(env.RabbitURL)
	if err != nil {
		log.Fatalf("failed to initialise rabbitmq: %s reason: %s\n", mm7RPC.Name, err)
	}

	redisCon, err := redis.Connect(env.RedisURL)
	if err != nil {
		log.Fatalf("failed to initialise redis: %s reason: %s\n", mm7RPC.Name, err)
	}

	limiter, err := redis.NewLimiter(env.RedisURL)
	if err != nil {
		log.Fatalf("failed to initialise redis: %s reason: %s\n", mm7RPC.Name, err)
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
		log.Fatalf("failed to initialise s3: %s reason: %s\n", mm7RPC.Name, err)
	}

	port := env.RPCPort

	rabbitOpts := mm7RPC.RabbitPublishOptions{
		Exchange:     env.RabbitExchange,
		ExchangeType: env.RabbitExchangeType,
	}

	configVar := mm7RPC.ConfigVar{
		AWSRegion:      env.AWSRegion,
		MMSMediaBucket: env.MMSMediaBucket,
	}

	server, err := rpc.NewServer(mm7RPC.NewService(rabbitmq, rabbitOpts, redisCon, limiter, s3Svc, configVar), port)
	if err != nil {
		log.Fatalf("failed to initialise service: %s reason: %s\n", mm7RPC.Name, err)
	}

	log.Printf("%s service initialised and available on port %d", mm7RPC.Name, port)
	server.Listen()
}
