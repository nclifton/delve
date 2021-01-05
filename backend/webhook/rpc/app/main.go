package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kelseyhightower/envconfig"

	"google.golang.org/grpc"

	"github.com/burstsms/mtmo-tp/backend/lib/jaeger"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/app/db"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/app/queue"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/app/service"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
)

type webhookEnv struct {
	RPCHost            string `envconfig:"RPC_HOST"`
	RPCPort            string `envconfig:"RPC_PORT"`
	RabbitURL          string `envconfig:"RABBIT_URL"`
	PostgresURL        string `envconfig:"POSTGRES_URL"`
	RabbitExchange     string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType string `envconfig:"RABBIT_EXCHANGE_TYPE"`
}

func main() {
	var env webhookEnv
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", env.RPCHost, env.RPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	tracer, closer, err := jaeger.Connect(env.RPCHost)
	if err != nil {
		log.Fatalf("failed to init jaeger: %s", err)
	}
	defer closer.Close()

	queueConn, err := rabbit.Connect(env.RabbitURL)
	if err != nil {
		log.Fatalf("failed to init rabbit: %s\n with error: %s", env.RabbitURL, err)
	}
	rqueue := queue.NewRabbitQueue(queueConn, rabbit.PublishOptions{
		Exchange:     env.RabbitExchange,
		ExchangeType: env.RabbitExchangeType,
		Tracer:       tracer,
	})

	sqlDB, err := pgxpool.Connect(context.Background(), env.PostgresURL)
	if err != nil {
		log.Fatalf("failed to init postgres: %s\n with error: %s", env.PostgresURL, err)
	}
	pdb := db.NewSQLDB(sqlDB)

	webhookpb.RegisterServiceServer(grpcServer, service.NewWebhookService(pdb, rqueue))

	log.Printf("running: %s\n on port: %s\n", env.RPCHost, env.RPCPort)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to start grpc server for service: %s\n on port: %s\n", env.RPCHost, env.RPCPort)
	}
}
