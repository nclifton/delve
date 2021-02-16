// +build integration

package test

import (
	"os"

	"testing"

	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/builder"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/post/postbuilder"

	"github.com/burstsms/mtmo-tp/backend/lib/fixtures"
	"github.com/burstsms/mtmo-tp/backend/lib/rpcbuilder"
	"github.com/burstsms/mtmo-tp/backend/lib/workerbuilder"
)

var tfx *fixtures.TestFixtures

func TestMain(m *testing.M) {
	setupFixtures()
	code := m.Run()
	defer os.Exit(code)
	defer tfx.Teardown()
}

func setupFixtures() {
	tfx = fixtures.New(fixtures.Config{Name: "webhook"})
	tfx.SetupPostgres("webhook")
	tfx.SetupRabbit()
	tfx.SetupRedis()
	tfx.GRPCStart(webhookRPCService())
	tfx.StartWorker(fixtures.WorkerConfig{
		ContainerName:  "webhook-post-worker",
		RabbitExchange: "webhook",
		QueueName:      "webhook",
	}, webhookPostService())
}

func webhookRPCService() rpcbuilder.Service {
	return builder.NewService(builder.Config{
		PostRabbitExchange:     "webhook",
		PostRabbitExchangeType: "direct",
	})
}

func webhookPostService() workerbuilder.Service {

	service := postbuilder.New(postbuilder.Config{
		ClientTimeout: 3,
		RedisURL:      tfx.Redis.Address,
	})

	return service
}
