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
	"github.com/burstsms/mtmo-tp/backend/lib/jaeger"
	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

type Config struct {
	ContainerName               string `envconfig:"CONTAINER_NAME"`
	ContainerPort               int    `envconfig:"CONTAINER_PORT"`
	RabbitURL                   string `envconfig:"RABBIT_URL"`
	PostgresURL                 string `envconfig:"POSTGRES_URL"`
	TracerDisable               bool   `envconfig:"TRACER_DISABLE"`
	RabbitIgnoreClosedQueueConn bool   `envconfig:"RABBIT_IGNORE_CLOSED_QUEUE_CONN"`
	// special env variable that should be empty under Kubernetes so that the RPC listener will attach to all available ip addresses on the server
	// in the docker-compose env (example dev) we have to specify the host address
	DevHost string `envconfig:"DEV_HOST"`
}

type Deps struct {
	Tracer       opentracing.Tracer
	RabbitConn   wabbit.Conn
	PostgresConn *pgxpool.Pool
	Server       *grpc.Server
}

type grpcServer struct {
	conf         Config
	log          *logger.StandardLogger
	tracer       opentracing.Tracer
	tracerCloser io.Closer
	rabbitConn   wabbit.Conn
	postgresConn *pgxpool.Pool

	lis        net.Listener
	server     *grpc.Server
	serverOpts []grpc.ServerOption

	service Service
}

type Service interface {
	Run(deps Deps) error
}

func NewGRPCServer(config Config, service Service) grpcServer {
	return grpcServer{conf: config,
		log:     logger.NewLogger(),
		service: service,
	}
}

func NewGRPCServerFromEnv(service Service) grpcServer {
	stLog := logger.NewLogger()

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		stLog.Fatalf(context.Background(), "NewGRPCServerFromEnv", "failed to read env vars: %s", err)
	}

	return grpcServer{conf: config,
		log:     stLog,
		service: service,
	}
}

func (g *grpcServer) TracerClose() error {
	err := g.tracerCloser.Close()
	if err != nil {
		return fmt.Errorf("failed to close jaeger conn: %s", err)
	}

	return nil
}

func (g *grpcServer) createJaegerConn(ctx context.Context) error {
	if g.conf.TracerDisable {
		return nil
	}

	g.log.Fields(ctx, logger.Fields{
		"host": g.conf.ContainerName,
		"port": g.conf.ContainerPort}).Infof("Starting tracer connection")

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

func (g *grpcServer) createPostgresConn(ctx context.Context) error {
	if g.conf.PostgresURL == "" {
		return nil
	}

	g.log.Fields(ctx, logger.Fields{
		"host": g.conf.ContainerName,
		"port": g.conf.ContainerPort}).Infof("Starting db connection")

	postgresConn, err := pgxpool.Connect(ctx, g.conf.PostgresURL)
	if err != nil {
		return fmt.Errorf("failed to init postgres: %s\n with error: %s", g.conf.PostgresURL, err)
	}

	g.postgresConn = postgresConn
	return nil
}

func (g *grpcServer) createRabbitConn(ctx context.Context) error {
	if g.conf.RabbitURL == "" {
		return nil
	}

	g.log.Fields(ctx, logger.Fields{
		"host": g.conf.ContainerName,
		"port": g.conf.ContainerPort}).Infof("Starting rabbit connection")

	rabbitConn, err := rabbit.Connect(g.conf.RabbitURL, g.conf.RabbitIgnoreClosedQueueConn)
	if err != nil {
		return fmt.Errorf("failed to init rabbit: %s\n with error: %s", g.conf.RabbitURL, err)
	}

	g.rabbitConn = rabbitConn
	return nil
}

func (g *grpcServer) SetCustomListener(lis net.Listener) {
	g.lis = lis
}

func (g *grpcServer) createListener(ctx context.Context) error {
	g.log.Fields(ctx, logger.Fields{
		"host": g.conf.ContainerName,
		"port": g.conf.ContainerPort}).Infof("Starting listener")

	// TODO: remove DevHost once we ditch docker-compose - listener should not be provided a host name
	// listen host (DevHost) can be an empty string (https://golang.org/pkg/net/)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", g.conf.DevHost, g.conf.ContainerPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	g.lis = lis
	return nil
}

func (g *grpcServer) Listener() net.Listener {
	return g.lis
}

func (g *grpcServer) createServer() {
	server := grpc.NewServer(g.serverOpts...)

	g.server = server
}

func (g *grpcServer) setupDeps(ctx context.Context) error {
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

func (g *grpcServer) Start() error {
	ctx := context.Background()

	if err := g.setupDeps(ctx); err != nil {
		return err
	}

	g.createServer()

	if err := g.service.Run(Deps{
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
		g.log.Fields(ctx, logger.Fields{
			"host": g.conf.ContainerName,
			"port": g.conf.ContainerPort}).Infof("Starting service")

		if err := g.server.Serve(g.lis); err != nil {
			g.log.Fields(ctx, logger.Fields{
				"host": g.conf.ContainerName,
				"port": g.conf.ContainerPort}).Fatalf("Failed to start grpc server")
		}
	}()

	sigint := make(chan os.Signal, 1)

	// wait for Control c or sigterm/sighup signal to exit
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	// Block until a signal is received
	<-sigint

	g.stop(ctx)

	return nil
}

func (g *grpcServer) stop(ctx context.Context) {
	logService := g.log.Fields(ctx, logger.Fields{
		"host": g.conf.ContainerName,
		"port": g.conf.ContainerPort})

	logService.Infof("Stopping service")

	g.server.GracefulStop()

	if g.lis != nil {
		logService.Infof("Closing listener")

		g.lis.Close()
	}

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

	logService.Infof("End of service")
}
