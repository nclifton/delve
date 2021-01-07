// +build integration

package test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
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
func startGrpcServer(tfx *fixtures.TestFixtures) *bufconn.Listener {
	env := getWebhookEnv()
	app := service.New()
	app.Env = &service.WebhookEnv{
		RPCHost:            "",
		RPCPort:            "",
		RabbitURL:          tfx.Rabbit.ConnStr,
		PostgresURL:        tfx.Postgres.ConnStr,
		RabbitExchange:     env.RabbitExchange,
		RabbitExchangeType: env.RabbitExchangeType,
	}
	listener := bufconn.Listen(1024 * 1024)
	app.Listener = listener

	// we need to use a rabbitmq connection that will not do a os.Exit() when we stop the service
	app.IgnoreClosedQueueConnection = true
	go app.Run()

	return listener
}

func newSetup(t *testing.T, tfx *fixtures.TestFixtures, listener *bufconn.Listener) *testDeps {

	setup := &testDeps{
		ctx:      context.Background(),
		tfx:      tfx,
		env:      getWebhookEnv(),
		listener: listener,
	}

	setup.adb = assertdb.New(t, setup.tfx.Postgres.ConnStr)

	return setup
}

func (setup *testDeps) teardown(t *testing.T) {
	if setup.adb != nil {
		setup.adb.Teardown()
	}
	if setup.httpServer != nil {
		setup.httpServer.Close()
	}

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
	wkr := worker.New()
	wkr.Env = &worker.WebhookEnv{
		RPCPort:         0,
		RPCHost:         "",
		RabbitURL:       setup.tfx.Rabbit.ConnStr,
		ClientTimeout:   3,
		WorkerQueueName: "webhook.post",
		RedisURL:        setup.tfx.Redis.Address,
		NRName:          "",
		NRLicense:       "",
		NRTracing:       false,
	}
	wkr.IgnoreClosedQueueConnection = true
	// use go routine to run the webhook worker
	go wkr.Run()
	time.Sleep(100 * time.Millisecond) // wait a bit for the worker to become ready

}

func (setup *testDeps) startHttpServer(t *testing.T) {
	setup.httpRequests = make([]*http.Request, 0)
	setup.httpRequestBodies = make([]string, 0)
	setup.httpServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setup.httpRequests = append(setup.httpRequests, r)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("failed to read request body, %+v", err)
		}
		setup.httpRequestBodies = append(setup.httpRequestBodies, string(body))
		fmt.Fprintln(w, "thank you")
	}))
}

type ExpectedRequests struct {
	numberOfRequests int
	waitMilliseconds int
	methods          []string
	contentTypes     []string
	bodies           []string
}

func (setup *testDeps) waitForRequests(t *testing.T, want ExpectedRequests) {
	var cnt = 0
	log.Println("waiting for http request")
	for len(setup.httpRequests) < want.numberOfRequests || want.numberOfRequests == 0 {
		if cnt > want.waitMilliseconds {
			plural := ""
			if want.numberOfRequests > 1 {
				plural = "s"
			}
			if want.numberOfRequests > len(setup.httpRequests) {
				assert.Fail(t, fmt.Sprintf("timed out waiting for %d request%s", want.numberOfRequests, plural))
			} else {
				log.Printf("http request wait timed out")
				return
			}
		}
		time.Sleep(time.Millisecond)
		cnt++
	}
	log.Printf("received http request after %d milliseconds", cnt)
}

func (setup *testDeps) resetHttpRequests(t *testing.T) {
	setup.httpRequests = nil
	setup.httpRequestBodies = nil
}

func (setup *testDeps) marshalJson(t *testing.T, v interface{}) string {
	expectedBody, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("error: %+v", err)
	}
	return string(expectedBody)
}

func (setup *testDeps) assertRequests(t *testing.T, want ExpectedRequests) {
	assert.Equal(t, want.numberOfRequests, len(setup.httpRequests), "number of requests received")
	for i, method := range want.methods {
		req := setup.httpRequests[i]
		assert.Equal(t, req.Method, method, fmt.Sprintf("request %d method", i+1))
	}
	for i, contentType := range want.contentTypes {
		req := setup.httpRequests[i]
		assert.Equal(t, req.Header.Get("Content-Type"), contentType, fmt.Sprintf("request %d Content-Type", i+1))
	}
	for i, body := range want.bodies {
		log.Println(setup.httpRequestBodies[i])
		assert.JSONEq(t, body, setup.httpRequestBodies[i], fmt.Sprintf("request %d body json", i+1))
	}
}
