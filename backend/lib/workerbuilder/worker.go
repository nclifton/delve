package workerbuilder

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/burstsms/mtmo-tp/backend/lib/jaeger"
	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/nr"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"

	"github.com/NeowayLabs/wabbit"
	"github.com/kelseyhightower/envconfig"
	"github.com/opentracing/opentracing-go"
)

type Config struct {
	WorkerName string `envconfig:"WORKER_NAME"`
	RabbitURL  string `envconfig:"RABBIT_URL"`

	TracerDisable               bool `envconfig:"TRACER_DISABLE"`
	RabbitIgnoreClosedQueueConn bool `envconfig:"RABBIT_IGNORE_CLOSED_QUEUE_CONN"`

	NRName    string `envconfig:"NR_NAME"`
	NRLicense string `envconfig:"NR_LICENSE"`
	NRTracing bool   `envconfig:"NR_TRACING"`
}

type Deps struct {
	Worker *rabbit.Worker
}

type worker struct {
	conf    Config
	log     *logger.StandardLogger
	service Service

	tracer       opentracing.Tracer
	tracerCloser io.Closer
	rabbitConn   wabbit.Conn
	worker       *rabbit.Worker
}

type Service interface {
	Run(deps Deps) error
}

func NewWorker(config Config, service Service) *worker {
	return &worker{
		conf:    config,
		log:     logger.NewLogger(),
		service: service,
	}
}

func NewWorkerFromEnv(service Service) worker {
	stLog := logger.NewLogger()

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		stLog.Fatalf(context.Background(), "envconfig.Process", "failed to read env vars: %s", err)
	}

	return worker{
		conf:    config,
		log:     stLog,
		service: service,
	}
}

func (w *worker) Start() error {
	ctx := context.Background()

	if err := w.setupDeps(ctx); err != nil {
		return err
	}

	w.createWorker(w.conf.WorkerName)

	go func() {
		if err := w.service.Run(Deps{
			Worker: w.worker,
		}); err != nil {
			w.log.Fields(ctx, logger.Fields{
				"worker": w.conf.WorkerName}).Fatalf("Failed to start worker")
		}
	}()

	sigint := make(chan os.Signal, 1)

	// wait for Control c or sigterm/sighup signal to exit
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	// Block until a signal is received
	<-sigint

	w.stop(ctx)

	return nil
}

func (w *worker) stop(ctx context.Context) {
	logService := w.log.Fields(ctx, logger.Fields{
		"worker": w.conf.WorkerName})

	logService.Infof("Stopping worker")

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

	logService.Infof("End of Worker")
}

func (w *worker) createWorker(queueName string) {
	nrOpts := &nr.Options{
		AppName:                  w.conf.NRName,
		NewRelicLicense:          w.conf.NRLicense,
		DistributedTracerEnabled: w.conf.NRTracing,
	}

	worker := rabbit.NewWorkerWithTracer(queueName, w.rabbitConn, nrOpts, w.tracer)

	w.worker = worker
}

func (w *worker) setupDeps(ctx context.Context) error {
	if err := w.createJaegerConn(ctx); err != nil {
		return err
	}

	return w.createRabbitConn(ctx)
}

func (w *worker) createJaegerConn(ctx context.Context) error {
	if w.conf.TracerDisable {
		return nil
	}

	w.log.Fields(ctx, logger.Fields{
		"worker": w.conf.WorkerName}).Infof("Starting tracer connection")

	tracer, closer, err := jaeger.Connect(w.conf.WorkerName)
	if err != nil {
		return fmt.Errorf("failed to init jaeger: %s", err)
	}

	w.tracer = tracer
	w.tracerCloser = closer

	return nil
}

func (w *worker) createRabbitConn(ctx context.Context) error {
	w.log.Fields(ctx, logger.Fields{
		"worker": w.conf.WorkerName}).Infof("Starting rabbit connection")

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
