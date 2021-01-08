package servicebuilder

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/NeowayLabs/wabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/jaeger"
	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

func NewGRPCServer(config Config) grpcServer {
	return grpcServer{conf: config}
}

func NewGRPCServerFromEnv() grpcServer {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal("failed to read env vars:", err)
	}
	return grpcServer{conf: config}
}

type Config struct {
	RPCHost     string `envconfig:"RPC_HOST"`
	RPCPort     string `envconfig:"RPC_PORT"`
	RabbitURL   string `envconfig:"RABBIT_URL"`
	PostgresURL string `envconfig:"POSTGRES_URL"`

	TracerDisable               bool `envconfig:"TRACER_DISABLE"`
	RabbitIgnoreClosedQueueConn bool `envconfig:"RABBIT_IGNORE_CLOSED_QUEUE_CONN"`
}

type Deps struct {
	Tracer       opentracing.Tracer
	RabbitConn   wabbit.Conn
	PostgresConn *pgxpool.Pool
	Server       *grpc.Server
}

type grpcServer struct {
	conf         Config
	tracer       opentracing.Tracer
	tracerCloser io.Closer
	rabbitConn   wabbit.Conn
	postgresConn *pgxpool.Pool

	lis        net.Listener
	server     *grpc.Server
	serverOpts []grpc.ServerOption
}

func (g *grpcServer) TracerClose() error {
	err := g.tracerCloser.Close()
	if err != nil {
		return fmt.Errorf("failed to close jaeger conn: %s", err)
	}

	return nil
}

func (g *grpcServer) createJaegerConn() error {
	if g.conf.TracerDisable {
		return nil
	}

	tracer, closer, err := jaeger.Connect(g.conf.RPCHost)
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

func (g *grpcServer) createPostgresConn() error {
	if g.conf.PostgresURL == "" {
		return nil
	}
	postgresConn, err := pgxpool.Connect(context.Background(), g.conf.PostgresURL)
	if err != nil {
		return fmt.Errorf("failed to init postgres: %s\n with error: %s", g.conf.PostgresURL, err)
	}

	g.postgresConn = postgresConn
	return nil
}

func (g *grpcServer) createRabbitConn() error {
	if g.conf.RabbitURL == "" {
		return nil
	}
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

func (g *grpcServer) createListener() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", g.conf.RPCHost, g.conf.RPCPort))
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

func (g *grpcServer) setupDeps() error {
	var err error

	err = g.createJaegerConn()
	if err != nil {
		return err
	}

	err = g.createPostgresConn()
	if err != nil {
		return err
	}

	err = g.createRabbitConn()
	if err != nil {
		return err
	}

	return nil
}

func (g *grpcServer) GRPCStart(registerCB func(deps Deps) error) error {
	err := g.setupDeps()
	if err != nil {
		return err
	}
	g.createServer()
	err = registerCB(Deps{
		Tracer:       g.tracer,
		RabbitConn:   g.rabbitConn,
		PostgresConn: g.postgresConn,
		Server:       g.server,
	})
	if err != nil {
		return err
	}
	if g.lis == nil {
		err = g.createListener()
		if err != nil {
			return err
		}
	}
	err = g.server.Serve(g.lis)
	if err != nil {
		return fmt.Errorf("failed to start grpc server for service: %s\n on port: %s\n", g.conf.RPCHost, g.conf.RPCPort)
	}

	return nil
}
