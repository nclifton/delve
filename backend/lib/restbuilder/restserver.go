package restbuilder

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/burstsms/mtmo-tp/backend/lib/jaeger"
	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/julienschmidt/httprouter"
	"github.com/kelseyhightower/envconfig"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func NewFromEnv(impl Impl) *restbuilder {
	var conf Config
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Fatal("Failed to read env vars:", err)
	}
	ctx := context.Background()
	return &restbuilder{impl: impl, conf: conf, ctx: ctx}
}

type Config struct {
	ContainerName string `envconfig:"CONTAINER_NAME"`
	ContainerPort string `envconfig:"CONTAINER_PORT"`
}

type Deps struct {
	Router *httprouter.Router
	Tracer opentracing.Tracer
	Log    *logger.StandardLogger
}

type Impl interface {
	Run(deps Deps) error
}

type restbuilder struct {
	conf         Config
	router       *httprouter.Router
	srv          *http.Server
	log          *logrus.Entry
	tracer       opentracing.Tracer
	tracerCloser io.Closer
	impl         Impl
	ctx          context.Context
}

func (r *restbuilder) Start() {
	r.log = logger.NewLogger().Fields(r.ctx, logger.Fields{
		"host": r.conf.ContainerName,
		"port": r.conf.ContainerPort})

	err := r.setup()
	if err != nil {
		r.log.Fatal(err)
	}

	appLog := logger.NewLogger()

	err = r.impl.Run(Deps{
		Router: r.router,
		Tracer: r.tracer,
		Log:    appLog,
	})
	if err != nil {
		r.log.Fatal(err)
	}

	r.serve()

	sigint := make(chan os.Signal, 1)

	// wait for Control c or sigterm/sighup signal to exit
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	// Block until a signal is received
	<-sigint

	r.stop()
}

func (r *restbuilder) setup() error {
	if r.tracer == nil {
		tracer, closer, err := jaeger.Connect(r.conf.ContainerName)
		if err != nil {
			return err
		}
		r.tracer = tracer
		r.tracerCloser = closer
	}

	r.router = httprouter.New()

	return nil
}

func (r *restbuilder) serve() {
	r.srv = &http.Server{
		Addr:    ":" + r.conf.ContainerPort,
		Handler: r.router,
	}

	r.log.Infof("Starting http server")
	go func() {
		if err := r.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			r.log.Fatalf("Failed to start http server: %s", err)
		}
	}()
}

func (r *restbuilder) stop() {
	r.log.Infof("Stopping http server")

	if err := r.srv.Shutdown(r.ctx); err != nil {
		r.log.Errorf("Failed to stop http server: %s", err)
	}

	r.log.Infof("Successfully stopped http server")

	if r.tracerCloser != nil {
		r.log.Infof("Closing tracer connection")

		if err := r.tracerCloser.Close(); err != nil {
			r.log.Errorf("Failed to close tracer connection: %s", err)
		}
	}

	r.log.Infof("End of service")
}
