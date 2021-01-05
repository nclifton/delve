package service

import (
	"context"
	"fmt"
	"io"
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
	env                         *WebhookEnv
	lis                         net.Listener
	tracer                      opentracing.Tracer
	tracerCloser                io.Closer
	queue                       queue.Queue
	ignoreClosedQueueConnection bool
	queueConn                   rabbit.Conn
	db                          db.DB
	sqlDB                       *pgxpool.Pool
	opts                        []grpc.ServerOption
	grpcServer                  *grpc.Server
}

func New() app {
	app := app{
		env:                         &WebhookEnv{},
		ignoreClosedQueueConnection: false,
	}
	app.loadEnv()
	return app
}

func (a *app) loadEnv() {
	if err := envconfig.Process("", a.env); err != nil {
		log.Fatal("failed to read env vars:", err)
	}
}

func (a *app) SetEnv(env *WebhookEnv) {
	a.env = env
}

func (a *app) checkDependencies() {
	if a.lis == nil {
		a.createListener()
	}
	if a.tracer == nil {
		a.createTracer()
	}
	if a.queue == nil {
		a.createQueue()
	}
	if a.db == nil {
		a.createDb()
	}
}

func (a *app) createListener() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", a.env.RPCHost, a.env.RPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	a.lis = lis
}

func (a *app) SetListener(listener net.Listener) {
	a.lis = listener
}

func (a *app) createTracer() {
	var err error
	a.tracer, a.tracerCloser, err = jaeger.Connect(a.env.RPCHost)
	if err != nil {
		log.Fatalf("failed to init jaeger: %s", err)
	}
}

func (a *app) SetTracer(tracer opentracing.Tracer) {
	a.tracer = tracer
}

func (a *app) createQueue() {
	var err error
	a.queueConn, err = rabbit.Connect(a.env.RabbitURL, a.ignoreClosedQueueConnection)
	if err != nil {
		log.Fatalf("failed to init rabbit: %s\n with error: %s", a.env.RabbitURL, err)
	}
	a.queue = queue.NewRabbitQueue(a.queueConn, rabbit.PublishOptions{
		Exchange:     a.env.RabbitExchange,
		ExchangeType: a.env.RabbitExchangeType,
		Tracer:       a.tracer,
	})
}

func (a *app) IgnoreClosedQueueConnection() {
	a.ignoreClosedQueueConnection = true
}

func (a *app) SetQueue(queue queue.Queue) {
	a.queue = queue
}

func (a *app) createDb() {
	var err error
	a.sqlDB, err = pgxpool.Connect(context.Background(), a.env.PostgresURL)
	if err != nil {
		log.Fatalf("failed to init postgres: %s\n with error: %s", a.env.PostgresURL, err)
	}
	a.db = db.NewSQLDB(a.sqlDB)
}

func (a *app) SetDb(db db.DB) {
	a.db = db
}

func (a *app) SetServerOpts(opts ...grpc.ServerOption) {
	a.opts = opts
}

func (a *app) Run() {
	a.prepareServer()
	if err := a.grpcServer.Serve(a.lis); err != nil {
		log.Fatalf("failed to start grpc server for service: %s\n on port: %s\n",
			a.env.RPCHost, a.env.RPCPort)
	}
	log.Printf("running: %s\n on port: %s\n", a.env.RPCHost, a.env.RPCPort)
}

func (a *app) prepareServer() {
	a.checkDependencies()
	a.grpcServer = grpc.NewServer(a.opts...)
	webhookpb.RegisterServiceServer(a.grpcServer, NewWebhookService(a.db, a.queue))
}
func (a *app) Start() {
	a.prepareServer()
	go func() {
		if err := a.grpcServer.Serve(a.lis); err != nil {
			log.Fatalf("failed to start grpc server for service: %s\n on port: %s\n",
				a.env.RPCHost, a.env.RPCPort)
		}
		log.Printf("running: %s on port: %s\n", a.env.RPCHost, a.env.RPCPort)
	}()
}

func (a *app) Close() error {

	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
	}

	// close connections being used by dependencies
	if a.sqlDB != nil {
		a.sqlDB.Close()
	}
	if a.queueConn != nil {
		a.queueConn.Close()
	}

	if a.tracerCloser != nil {
		err := a.tracerCloser.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
