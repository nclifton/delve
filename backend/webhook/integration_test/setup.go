// +build integration

package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/burstsms/mtmo-tp/backend/webhook/integration_test/assertdb"
	"github.com/burstsms/mtmo-tp/backend/webhook/integration_test/fixtures"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/app/service"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker"
)

var testStartTime time.Time

func init() {
	testStartTime = time.Now()
}

type testDeps struct {
	ctx               context.Context
	tfx               *fixtures.TestFixtures
	env               *WebhookEnv
	listener          *bufconn.Listener
	connectionToRPC   *grpc.ClientConn
	adb               *assertdb.AssertDb
	appClose          func() error
	httpServer        *httptest.Server
	httpRequests      []*http.Request
	httpRequestBodies []string
}

type WebhookEnv struct {
	RabbitExchange     string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType string `envconfig:"RABBIT_EXCHANGE_TYPE"`
	MigrationRoot      string `envconfig:"MIGRATION_ROOT"`
}

func getWebhookEnv() *WebhookEnv {
	var env WebhookEnv
	if err := envconfig.Process("webhook", &env); err != nil {
		log.Fatal("failed to read env vars:", err)
	}
	return &env
}

func (setup *testDeps) teardown(t *testing.T) {
	if setup.adb != nil {
		setup.adb.Teardown()
	}
	setup.connectionToRPC.Close()
	if setup.appClose != nil {
		err := setup.appClose()
		if err != nil {
			t.Fatalf("failed to close application: %+v", err)
		}
	}
	if setup.httpServer != nil {
		setup.httpServer.Close()
	}

}

func setupForTest(t *testing.T, tfx *fixtures.TestFixtures) *testDeps {

	setup := &testDeps{
		ctx: context.Background(),
		tfx: tfx,
		env: getWebhookEnv(),
	}

	app := service.New()
	app.SetEnv(&service.WebhookEnv{
		RPCHost:            "webhook service under test",
		RPCPort:            "N/A",
		RabbitURL:          tfx.Rabbit.ConnStr,
		PostgresURL:        tfx.Postgres.ConnStr,
		RabbitExchange:     setup.env.RabbitExchange,
		RabbitExchangeType: setup.env.RabbitExchangeType,
	})
	app.SetListener(setup.getBufListener())
	app.SetTracer(setup.getNoopTracer())
	// we need to use a rabbitmq connection that will not fatal the test when we stop the service
	app.IgnoreClosedQueueConnection()
	go app.Run()
	setup.appClose = app.Close

	setup.adb = assertdb.New(t, setup.tfx.Postgres.ConnStr)

	return setup
}

func (setup *testDeps) getNoopTracer() opentracing.Tracer {
	return opentracing.GlobalTracer()
}

func (setup *testDeps) getBufListener() *bufconn.Listener {
	bufferSize := 1024 * 1024
	setup.listener = bufconn.Listen(bufferSize)
	return setup.listener
}

func (setup *testDeps) getClient(t *testing.T) webhookpb.ServiceClient {
	conn, err := grpc.DialContext(setup.ctx, "",
		grpc.WithContextDialer(getBufDialer(setup.listener)),
		grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial: %+v", err)
	}
	setup.connectionToRPC = conn
	return webhookpb.NewServiceClient(setup.connectionToRPC)

}

func getBufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}

func (setup *testDeps) startWorker(t *testing.T) {
	// use go routine to start the webhook worker
	wkr := worker.New()
	wkr.SetEnv(&worker.WebhookEnv{
		RPCPort:         0,
		RPCHost:         "n/a",
		RabbitURL:       setup.tfx.Rabbit.ConnStr,
		ClientTimeout:   3,
		WorkerQueueName: "webhook.post",
		RedisURL:        setup.tfx.Redis.Address,
		NRName:          "",
		NRLicense:       "",
		NRTracing:       false,
	})
	wkr.IgnoreClosedQueueConnection()
	go wkr.Run()

}

func (setup *testDeps) startHttpServer(t *testing.T) {
	setup.httpRequests = make([]*http.Request, 0)
	setup.httpRequestBodies = make([]string, 0)
	setup.httpServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setup.httpRequests = append(setup.httpRequests, r)
		body,err := ioutil.ReadAll(r.Body)
		if err != nil{
			log.Printf("failed to read request body, %+v", err)
		}
		setup.httpRequestBodies = append(setup.httpRequestBodies,string(body))
		fmt.Fprintln(w, "Hello, client")
	}))
}
