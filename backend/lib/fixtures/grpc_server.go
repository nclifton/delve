package fixtures

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/burstsms/mtmo-tp/backend/lib/rpcbuilder"
)

func (tfx *TestFixtures) GRPCStart(service rpcbuilder.Service) {

	ctx := context.Background()
	hcPort := port(tfx.env.RPCHealthCheckPort)

	s := rpcbuilder.NewGRPCServer(
		ctx,
		rpcbuilder.Config{
			ContainerName:               "integration-test-service",
			ContainerPort:               0,
			RabbitURL:                   tfx.Rabbit.ConnStr,
			PostgresURL:                 tfx.Postgres.ConnStr,
			TracerDisable:               true,
			RabbitIgnoreClosedQueueConn: true,
			HealthCheckPort:             hcPort,
			MaxGoRoutines:               200,
			DevHost:                     tfx.env.HealthCheckHost,
		},
		service,
	)
	tfx.RPCHealthCheckURI = fmt.Sprintf("http://%s:%s", tfx.env.HealthCheckHost, hcPort)

	s.SetCustomListener(bufconn.Listen(1024 * 1024))
	tfx.GRPCListener = s.Listener()
	tfx.teardowns = append(tfx.teardowns, func() {
		s.Stop(ctx)
		tfx.GRPCListener = nil
	})

	// start in a go routine
	go func() {
		err := s.Start(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

}

func (tfx *TestFixtures) GRPCClientConnection(t *testing.T, ctx context.Context) *grpc.ClientConn {
	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(getBufDialer(tfx.GRPCListener.(*bufconn.Listener))),
		grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial: %+v", err)
	}
	return conn
}

func getBufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}
