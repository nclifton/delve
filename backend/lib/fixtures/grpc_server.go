package fixtures

import (
	"context"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/burstsms/mtmo-tp/backend/lib/servicebuilder"
)

func (tfx *TestFixtures) GRPCStart(registerCB func(deps servicebuilder.Deps) error) {

	s := servicebuilder.NewGRPCServer(servicebuilder.Config{
		RPCHost:                     "",
		RPCPort:                     "",
		RabbitURL:                   tfx.Rabbit.ConnStr,
		PostgresURL:                 tfx.Postgres.ConnStr,
		TracerDisable:               true,
		RabbitIgnoreClosedQueueConn: true})

	s.SetCustomListener(bufconn.Listen(1024 * 1024))
	tfx.GRPCListener = s.Listener()

	// start in a go routine
	go func() {
		err := s.GRPCStart(registerCB)
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
