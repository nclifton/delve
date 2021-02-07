package workerbuilder

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/burstsms/mtmo-tp/backend/lib/health"
	"github.com/burstsms/mtmo-tp/backend/lib/jaeger"
	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"

	"github.com/NeowayLabs/wabbit"
	"github.com/kelseyhightower/envconfig"
	"github.com/opentracing/opentracing-go"
)

type Config struct {
	ContainerName               string `envconfig:"CONTAINER_NAME"`
	RabbitURL                   string `envconfig:"RABBIT_URL"`
	TracerDisable               bool   `envconfig:"TRACER_DISABLE"`
	RabbitIgnoreClosedQueueConn bool   `envconfig:"RABBIT_IGNORE_CLOSED_QUEUE_CONN"`
	NRName                      string `envconfig:"NR_NAME"`
	NRLicense                   string `envconfig:"NR_LICENSE"`
	NRTracing                   bool   `envconfig:"NR_TRACING"`
	RabbitQueueName             string `envconfig:"RABBIT_QUEUE_NAME"`
	RabbitExchange              string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType          string `envconfig:"RABBIT_EXCHANGE_TYPE"`
	RabbitPrefetchedCount       int    `envconfig:"RABBIT_PREFETCHED_COUNT"`
	HealthCheckPort             string `envconfig:"HEALTH_CHECK_PORT" default:"8086"`
	MaxGoRoutines               int    `envconfig:"MAX_GO_ROUTINES" default:"100"`
	// special env variable that should be empty under Kubernetes so that the Health Check listener will attach to all available ip addresses on the server
	HealthCheckHost string `envconfig:"HEALTH_CHECK_HOST" default:""`
}

type Deps struct {
	Worker         *rabbit.Worker
	ConsumeOptions rabbit.ConsumeOptions
	Health         health.HealthCheckService
}

type worker struct {
	conf         Config
	log          *logger.StandardLogger
	service      Service
	tracer       opentracing.Tracer
	tracerCloser io.Closer
	rabbitConn   wabbit.Conn
	worker       *rabbit.Worker
	health       health.HealthCheckService
}

type Service interface {
	Run(deps Deps) error
}

func NewWorker(ctx context.Context, config Config, service Service) *worker {
	// TODO consider allowing injection of logger and health check service through this or another constructor

	stLog := logger.NewLogger()
	return &worker{
		conf:    config,
		log:     stLog,
		service: service,
		health:  health.New(ctx, healthCheckConfig(config), stLog),
	}
}

func NewWorkerFromEnv(ctx context.Context, service Service) worker {
	stLog := logger.NewLogger()

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		stLog.Fatalf(context.Background(), "envconfig.Process", "failed to read env vars: %s", err)
	}

	return worker{
		conf:    config,
		log:     stLog,
		service: service,
		health:  health.New(ctx, healthCheckConfig(config), stLog),
	}
}

func healthCheckConfig(config Config) health.Config {
	return health.Config{
		Host:          config.HealthCheckHost,
		Port:          config.HealthCheckPort,
		MaxGoRoutines: config.MaxGoRoutines,
		LoggerFields:  loggerFields(config),
	}
}

func loggerFields(config Config) logger.Fields {
	lf := logger.Fields{
		"worker": config.ContainerName}
	return lf
}

func (w *worker) Start(ctx context.Context) error {

	if err := w.createJaegerConn(ctx); err != nil {
		return err
	}

	if err := w.createRabbitConn(ctx); err != nil {
		return err
	}

	w.createWorker()

	go func() {
		if err := w.service.Run(Deps{
			Worker: w.worker,
			ConsumeOptions: rabbit.ConsumeOptions{
				PrefetchCount:        w.conf.RabbitPrefetchedCount,
				Exchange:             w.conf.RabbitExchange,
				ExchangeType:         w.conf.RabbitExchangeType,
				QueueName:            w.conf.RabbitQueueName,
				RetryScale:           rabbit.RetryScale,
				AllowConnectionClose: w.conf.RabbitIgnoreClosedQueueConn,
			},
			Health: w.health,
		}); err != nil {
			w.log.Fields(ctx, loggerFields(w.conf)).Fatalf("Failed to start worker")
		}
	}()

	sigint := make(chan os.Signal, 1)

	// wait for Control c or sigterm/sighup signal to exit
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	// Block until a signal is received
	<-sigint

	w.Stop(ctx)

	return nil
}

func (w *worker) Stop(ctx context.Context) {

	w.health.SetServiceReady(false)

	logService := logFields(w, ctx)

	logService.Infof("Stopping worker")

	//TODO we need a graceful way to stop the consumer
	// just closing the connection could cause the consumer to do an os.Exit which is not considered playing nice

	if w.rabbitConn != nil {
		logService.Infof("Closing rabbit connection")
		w.rabbitConn.Close()
	}

	if w.tracerCloser != nil {
		logService.Infof("Closing tracer connection")
		if err := w.TracerClose(); err != nil {
			logService.Fatalf("Failed to close tracer connection")
		}
	}

	w.health.Stop(ctx)

	logService.Infof("End of Worker")
}

func logFields(w *worker, ctx context.Context) *logrus.Entry {
	logService := w.log.Fields(ctx, loggerFields(w.conf))
	return logService
}

func (w *worker) createWorker() {
	nrOpts := &nr.Options{
		AppName:                  w.conf.NRName,
		NewRelicLicense:          w.conf.NRLicense,
		DistributedTracerEnabled: w.conf.NRTracing,
	}

	w.worker = rabbit.NewWorkerWithTracer(w.conf.ContainerName, w.rabbitConn, nrOpts, w.tracer)
}

func (w *worker) createJaegerConn(ctx context.Context) error {
	if w.conf.TracerDisable {
		return nil
	}

	w.log.Fields(ctx, loggerFields(w.conf)).Infof("Starting tracer connection: '%s'", w.conf.ContainerName)

	tracer, closer, err := jaeger.Connect(w.conf.ContainerName)
	if err != nil {
		return fmt.Errorf("failed to init jaeger: %s", err)
	}

	w.tracer = tracer
	w.tracerCloser = closer

	return nil
}

func (w *worker) createRabbitConn(ctx context.Context) error {
	w.log.Fields(ctx, logger.Fields{
		"worker": w.conf.ContainerName}).Infof("Starting rabbit connection")

	rabbitConn, err := rabbit.Connect(w.conf.RabbitURL, w.conf.RabbitIgnoreClosedQueueConn)
	if err != nil {
		return fmt.Errorf("failed to init rabbit: %s\n with error: %s", w.conf.RabbitURL, err)
	}

	w.rabbitConn = rabbitConn
	return nil
}

func (w *worker) TracerClose() error {
	err := w.tracerCloser.Close()
	if err != nil {
		return fmt.Errorf("failed to close jaeger conn: %s", err)
	}

	return nil
}
