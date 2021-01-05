package service

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kelseyhightower/envconfig"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"

	"github.com/burstsms/mtmo-tp/backend/lib/jaeger"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/app/db"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/app/queue"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
)

type webhookImpl struct {
	db    db.DB
	queue queue.Queue
	webhookpb.UnimplementedServiceServer
}

func NewWebhookService(db db.DB, queue queue.Queue) webhookpb.ServiceServer {
	return &webhookImpl{db: db, queue: queue}
}

type WebhookEnv struct {
	RPCHost            string `envconfig:"RPC_HOST"`
	RPCPort            string `envconfig:"RPC_PORT"`
	RabbitURL          string `envconfig:"RABBIT_URL"`
	PostgresURL        string `envconfig:"POSTGRES_URL"`
	RabbitExchange     string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType string `envconfig:"RABBIT_EXCHANGE_TYPE"`
}

type app struct {
	Env                         *WebhookEnv
	Listener                    net.Listener
	Tracer                      opentracing.Tracer
	IgnoreClosedQueueConnection bool
	Opts                        []grpc.ServerOption
}

func New() app {
	app := app{
		Env:                         &WebhookEnv{},
		IgnoreClosedQueueConnection: false,
	}
	if err := envconfig.Process("", app.Env); err != nil {
		log.Fatal("failed to read env vars:", err)
	}
	return app
}

func (a *app) Run() {
	var err error
	if a.Listener == nil {
		a.Listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", a.Env.RPCHost, a.Env.RPCPort))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
	}
	if a.Env.RPCHost != "" {
		a.Tracer, _, err = jaeger.Connect(a.Env.RPCHost)
		if err != nil {
			log.Fatalf("failed to init jaeger: %s", err)
		}
		opentracing.SetGlobalTracer(a.Tracer)
	}
	queueConn, err := rabbit.Connect(a.Env.RabbitURL, a.IgnoreClosedQueueConnection)
	if err != nil {
		log.Fatalf("failed to init rabbit: %s\n with error: %s", a.Env.RabbitURL, err)
	}
	queue := queue.NewRabbitQueue(queueConn, rabbit.PublishOptions{
		Exchange:     a.Env.RabbitExchange,
		ExchangeType: a.Env.RabbitExchangeType,
		Tracer:       a.Tracer,
	})
	sqlDB, err := pgxpool.Connect(context.Background(), a.Env.PostgresURL)
	if err != nil {
		log.Fatalf("failed to init postgres: %s\n with error: %s", a.Env.PostgresURL, err)
	}
	grpcServer := grpc.NewServer(a.Opts...)
	webhookpb.RegisterServiceServer(grpcServer, NewWebhookService(db.NewSQLDB(sqlDB), queue))

	if err := grpcServer.Serve(a.Listener); err != nil {
		log.Fatalf("failed to start grpc server for service: %s\n on port: %s\n",
			a.Env.RPCHost, a.Env.RPCPort)
	}
	log.Printf("running: %s\n on port: %s\n", a.Env.RPCHost, a.Env.RPCPort)
}
