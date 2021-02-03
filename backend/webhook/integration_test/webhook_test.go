// +build integration

package test

import (
	"log"
	"os"

	"testing"

	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/builder"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/post/postbuilder"

	"github.com/burstsms/mtmo-tp/backend/lib/fixtures"
	"github.com/burstsms/mtmo-tp/backend/lib/redis"
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
	tfx = fixtures.New("webhook")
	tfx.SetupPostgres("webhook")
	tfx.SetupRabbit()
	tfx.SetupRedis()
	tfx.GRPCStart(webhookRPCService())
	tfx.StartWorker("webhook-post-worker-service", webhookPostService())
}

func webhookRPCService() rpcbuilder.Service {
	return builder.NewService(builder.Config{
		PostRabbitExchange:     "webhook",
		PostRabbitExchangeType: "direct",
	})
}

func webhookPostService() workerbuilder.Service {

	service := postbuilder.New(postbuilder.Config{
		ClientTimeout:         3,
		RedisURL:              tfx.Redis.Address,
		RabbitExchange:        "webhook",
		RabbitExchangeType:    "direct",
		RabbitPrefetchedCount: 1,
	})

	limiter, err := redis.NewLimiter(tfx.Redis.Address)
	if err != nil {
		log.Fatal(err)
	}
	service.SetLimiter(limiter)

	return service
}
