package rpcbuilder

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/NeowayLabs/wabbit"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/burstsms/mtmo-tp/backend/lib/health"
	"github.com/burstsms/mtmo-tp/backend/lib/jaeger"
	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
)

type Config struct {
	ContainerName               string `envconfig:"CONTAINER_NAME"`
	ContainerPort               int    `envconfig:"CONTAINER_PORT"`
	RabbitURL                   string `envconfig:"RABBIT_URL"`
	PostgresURL                 string `envconfig:"POSTGRES_URL"`
	TracerDisable               bool   `envconfig:"TRACER_DISABLE"`
	RabbitIgnoreClosedQueueConn bool   `envconfig:"RABBIT_IGNORE_CLOSED_QUEUE_CONN"`
	HealthCheckPort             string `envconfig:"HEALTH_CHECK_PORT" default:"8086"`
	MaxGoRoutines               int    `envconfig:"MAX_GO_ROUTINES" default:"100"`

	// special env variable that should be empty under Kubernetes so that the RPC and Health Check listeners will attach to all available ip addresses on the server
	// in the docker-compose env (example dev) we must specify the host address
	// this will also become the health check listen/serve host, if blank, health check host will be ContainerName
	DevHost string `envconfig:"DEV_HOST"`
}

type Deps struct {
	Health       health.HealthCheckService
	Tracer       opentracing.Tracer
	RabbitConn   wabbit.Conn
	PostgresConn *pgxpool.Pool
	Server       *grpc.Server
}

type rpcServerProperties struct {
	conf         Config
	health       health.HealthCheckService
	log          *logger.StandardLogger
	tracer       opentracing.Tracer
	tracerCloser io.Closer
	rabbitConn   wabbit.Conn
	postgresConn *pgxpool.Pool
	lis          net.Listener
	server       *grpc.Server
	serverOpts   []grpc.ServerOption
	service      Service
}

type Service interface {
	Run(deps Deps) error
}

func NewGRPCServer(ctx context.Context, config Config, service Service) rpcServerProperties {
	// TODO consider allowing injection of logger and health check service through this or another constructor
	stLog := logger.NewLogger()
	return rpcServerProperties{
		conf:    config,
		log:     stLog,
		service: service,
		health:  health.New(ctx, healthCheckConfig(config), stLog),
	}
}

func NewGRPCServerFromEnv(ctx context.Context, service Service) rpcServerProperties {
	stLog := logger.NewLogger()
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		stLog.Fatalf(ctx, "NewGRPCServerFromEnv", "failed to read env vars: %s", err)
	}
	return rpcServerProperties{
		conf:    config,
		log:     stLog,
		service: service,
		health:  health.New(ctx, healthCheckConfig(config), stLog),
	}
}

func healthCheckConfig(config Config) health.Config {
	return health.Config{
		Host: func() string {
			if config.DevHost == "" {
				return config.ContainerName
			}
			return config.DevHost
		}(),
		Port:          config.HealthCheckPort,
		MaxGoRoutines: config.MaxGoRoutines,
		LoggerFields:  loggerFields(config),
	}
}

func loggerFields(config Config) logger.Fields {
	return logger.Fields{
		"host": config.ContainerName,
		"port": config.ContainerPort}
}

func (g *rpcServerProperties) TracerClose() error {
	err := g.tracerCloser.Close()
	if err != nil {
		return fmt.Errorf("failed to close jaeger conn: %s", err)
	}

	return nil
}

func (g *rpcServerProperties) logFields(ctx context.Context) *logrus.Entry {
	return g.log.Fields(ctx, loggerFields(g.conf))
}

func (g *rpcServerProperties) createJaegerConn(ctx context.Context) error {
	if g.conf.TracerDisable {
		return nil
	}

	g.logFields(ctx).Infof("Starting tracer connection: %s", g.conf.ContainerName)

	tracer, closer, err := jaeger.Connect(g.conf.ContainerName)
	if err != nil {
		return fmt.Errorf("failed to init jaeger: %s", err)
	}

	g.tracer = tracer
	g.tracerCloser = closer
	g.serverOpts = append(g.serverOpts, grpc.UnaryInterceptor(
		otgrpc.OpenTracingServerInterceptor(g.tracer, otgrpc.LogPayloads()),
	))
	return nil
}

func (g *rpcServerProperties) createPostgresConn(ctx context.Context) error {
	if g.conf.PostgresURL == "" {
		return nil
	}

	// seeing the service uses postgres, let's add the health check for it.
	//  - will fail until it gets a working connection
	g.health.AddPostgresReadinessCheck()

	g.logFields(ctx).Infof("Starting db connection")

	postgresConn, err := pgxpool.Connect(ctx, g.conf.PostgresURL)
	if err != nil {
		return fmt.Errorf("failed to init postgres: %s\n with error: %s", g.conf.PostgresURL, err)
	}

	g.postgresConn = postgresConn

	// seeing the service uses postgres, let's add the health check for it
	g.health.AddPostgresReadinessCheckConnection(postgresConn)

	return nil
}

func (g *rpcServerProperties) createRabbitConn(ctx context.Context) error {
	if g.conf.RabbitURL == "" {
		return nil
	}

	g.logFields(ctx).Infof("Starting rabbit connection")

	rabbitConn, err := rabbit.Connect(g.conf.RabbitURL, g.conf.RabbitIgnoreClosedQueueConn)
	if err != nil {
		return fmt.Errorf("failed to init rabbit: %s\n with error: %s", g.conf.RabbitURL, err)
	}

	g.rabbitConn = rabbitConn

	return nil
}

func (g *rpcServerProperties) SetCustomListener(lis net.Listener) {
	g.lis = lis
}

func (g *rpcServerProperties) createListener(ctx context.Context) error {
	g.logFields(ctx).Infof("Starting listener")

	// TODO: remove DevHost once we ditch docker-compose - listener should not be provided a host name
	// listen host (DevHost) can be an empty string (https://golang.org/pkg/net/)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", g.conf.DevHost, g.conf.ContainerPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	g.lis = lis
	return nil
}

func (g *rpcServerProperties) Listener() net.Listener {
	return g.lis
}

func (g *rpcServerProperties) setupDeps(ctx context.Context) error {
	var err error

	err = g.createJaegerConn(ctx)
	if err != nil {
		return err
	}

	err = g.createPostgresConn(ctx)
	if err != nil {
		return err
	}

	err = g.createRabbitConn(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (g *rpcServerProperties) Start(ctx context.Context) error {

	if err := g.setupDeps(ctx); err != nil {
		return err
	}

	g.server = grpc.NewServer(g.serverOpts...)

	if err := g.service.Run(Deps{
		Health:       g.health,
		Tracer:       g.tracer,
		RabbitConn:   g.rabbitConn,
		PostgresConn: g.postgresConn,
		Server:       g.server,
	}); err != nil {
		return err
	}

	if g.lis == nil {
		if err := g.createListener(ctx); err != nil {
			return err
		}
	}

	go func() {
		g.logFields(ctx).Infof("Starting service")
		g.health.SetServiceReady(true)
		if err := g.server.Serve(g.lis); err != nil {
			g.logFields(ctx).Fatalf("Failed to start grpc server: %+v", err)
		}
	}()

	sigint := make(chan os.Signal, 1)

	// wait for Control c or sigterm/sighup signal to exit
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	// Block until a signal is received
	<-sigint

	g.Stop(ctx)

	return nil
}

func (g *rpcServerProperties) Stop(ctx context.Context) {
	logService := g.logFields(ctx)

	logService.Infof("Stopping service")

	g.server.GracefulStop() // this stops the listener as well

	if g.postgresConn != nil {
		logService.Infof("Closing db connection")
		g.postgresConn.Close()
	}

	if g.rabbitConn != nil {
		logService.Infof("Closing rabbit connection")
		g.rabbitConn.Close()
	}

	if g.tracerCloser != nil {
		logService.Infof("Closing tracer connection")
		if err := g.TracerClose(); err != nil {
			logService.Fatalf("Failed to close tracer connection")
		}
	}

	g.health.Stop(ctx)

	logService.Infof("End of service")
}
