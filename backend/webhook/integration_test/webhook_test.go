// +build integration

package test

import (
	"log"
	"os"
	"reflect"

	"testing"

	"gotest.tools/assert"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/builder"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
	"github.com/burstsms/mtmo-tp/backend/webhook/worker/post/postbuilder"

	"github.com/burstsms/mtmo-tp/backend/lib/assertdb"
	"github.com/burstsms/mtmo-tp/backend/lib/fixtures"
	"github.com/burstsms/mtmo-tp/backend/lib/redis"
	"github.com/burstsms/mtmo-tp/backend/lib/rpcbuilder"
	"github.com/burstsms/mtmo-tp/backend/lib/workerbuilder"
)

var tfx *fixtures.TestFixtures

func TestMain(m *testing.M) {
	tfx = fixtures.New()
	tfx.SetupPostgres("webhook")
	tfx.SetupRabbit()
	tfx.SetupRedis()
	tfx.GRPCStart(webhookRPCService())
	tfx.StartWorker("webhook-post-worker-service", webhookPostService())
	code := m.Run()
	defer os.Exit(code)
	defer tfx.Teardown()
}

func webhookRPCService() rpcbuilder.Service {
	return builder.NewBuilder(builder.Config{
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

func setupForInsert(t *testing.T) *testDeps {
	return newSetup(t, tfx)
}

func Test_Insert(t *testing.T) {
	i := setupForInsert(t)
	log.Println("test Insert")
	defer i.teardown(t)
	client := i.getClient(t)

	type wantErr struct {
		status *status.Status
		ok     bool
	}
	tests := []struct {
		name    string
		params  *webhookpb.InsertParams
		want    func(*webhookpb.InsertReply)
		wantErr wantErr
	}{
		{
			name: "happy",
			params: &webhookpb.InsertParams{
				AccountId: "1",
				Event:     "event",
				Name:      "name",
				URL:       "url",
				RateLimit: 1,
			},
			want: func(response *webhookpb.InsertReply) {

				assert.Equal(t, response.Webhook.GetAccountId(), "1", "AccountId")
				assert.Equal(t, response.Webhook.GetEvent(), "event", "Event")
				assert.Equal(t, response.Webhook.GetName(), "name", "Name")
				assert.Equal(t, response.Webhook.GetURL(), "url", "URL")
				assert.Equal(t, response.Webhook.GetRateLimit(), int32(1), "RateLimit")
				assert.Check(t, response.Webhook.GetCreatedAt().AsTime().After(testStartTime), "CreatedAt")
				assert.Check(t, response.Webhook.GetUpdatedAt().AsTime().After(testStartTime), "UpdatedAt")

				i.SeeInDatabase("webhook", assertdb.Criteria{
					"account_id":   "1",
					"event":        "event",
					"name":         "name",
					"url":          "url",
					"created_at >": testStartTime.Format(assertdb.SQLTimestampWithoutTimeZone),
					"updated_at >": testStartTime.Format(assertdb.SQLTimestampWithoutTimeZone),
				})

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.Insert(i.ctx, tt.params)
			if tt.wantErr.status != nil && err != nil {
				errStatus, ok := status.FromError(err)
				assert.Equal(t, ok, tt.wantErr.ok, "grpc ok")
				assert.Check(t, reflect.DeepEqual(tt.wantErr.status, errStatus), "grpc status")
			} else if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			tt.want(got)
		})
	}
}

func setupForFind(t *testing.T) *testDeps {
	setup := newSetup(t, tfx)

	setup.HaveInDatabase("webhook", assertdb.Row{
		"id":         32767,
		"account_id": "42",
		"event":      "event1",
		"name":       "name1",
		"url":        "url1",
		"rate_limit": 2,
		"created_at": "2021-01-12 22:41:42",
		"updated_at": "2021-01-13 22:25:25"})

	setup.HaveInDatabase("webhook", assertdb.Row{
		"id":         32768,
		"account_id": "42",
		"event":      "event",
		"name":       "name",
		"url":        "url",
		"rate_limit": 1,
		"created_at": "2021-01-12 22:42:42",
		"updated_at": "2021-01-13 22:24:24"})

	return setup
}

func Test_Find(t *testing.T) {
	log.Println("test Find")

	i := setupForFind(t)
	defer i.teardown(t)
	client := i.getClient(t)

	type wantErr struct {
		status *status.Status
		ok     bool
	}
	tests := []struct {
		name    string
		params  *webhookpb.FindParams
		want    func(*webhookpb.FindReply)
		wantErr wantErr
	}{
		{
			name: "happy find 2",
			params: &webhookpb.FindParams{
				AccountId: "42",
			},
			want: func(response *webhookpb.FindReply) {
				assert.Equal(t, len(response.GetWebhooks()), 2, "number of webhooks")
				assert.Equal(t, response.Webhooks[0].GetId(), int64(32767), "Id")
				assert.Equal(t, response.Webhooks[0].GetAccountId(), "42", "AccountId")
				assert.Equal(t, response.Webhooks[0].GetName(), "name1", "Name")
				assert.Equal(t, response.Webhooks[0].GetEvent(), "event1", "Event")
				assert.Equal(t, response.Webhooks[0].GetURL(), "url1", "URL")
				assert.Equal(t, response.Webhooks[0].GetRateLimit(), int32(2), "RateLimit")
				assert.Equal(t, response.Webhooks[0].GetCreatedAt().AsTime().Format(assertdb.SQLTimestampWithoutTimeZone), "2021-01-12 22:41:42", "CreatedAt")
				assert.Equal(t, response.Webhooks[0].GetUpdatedAt().AsTime().Format(assertdb.SQLTimestampWithoutTimeZone), "2021-01-13 22:25:25", "UpdatedAt")

				assert.Equal(t, response.Webhooks[1].GetId(), int64(32768), "Id")
				assert.Equal(t, response.Webhooks[1].GetAccountId(), "42", "AccountId")
				assert.Equal(t, response.Webhooks[1].GetName(), "name", "Name")
				assert.Equal(t, response.Webhooks[1].GetEvent(), "event", "Event")
				assert.Equal(t, response.Webhooks[1].GetURL(), "url", "URL")
				assert.Equal(t, response.Webhooks[1].GetRateLimit(), int32(1), "RateLimit")
				assert.Equal(t, response.Webhooks[1].GetCreatedAt().AsTime().Format(assertdb.SQLTimestampWithoutTimeZone), "2021-01-12 22:42:42", "CreatedAt")
				assert.Equal(t, response.Webhooks[1].GetUpdatedAt().AsTime().Format(assertdb.SQLTimestampWithoutTimeZone), "2021-01-13 22:24:24", "UpdatedAt")
			},
		},
		{
			name: "happy find none",
			params: &webhookpb.FindParams{
				AccountId: "4422",
			},
			want: func(response *webhookpb.FindReply) {
				assert.Equal(t, len(response.GetWebhooks()), 0, "number of webhooks")
			},
			wantErr: wantErr{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.Find(i.ctx, tt.params)
			if tt.wantErr.status != nil && err != nil {
				errStatus, ok := status.FromError(err)
				assert.Equal(t, ok, tt.wantErr.ok, "grpc ok")
				assert.Check(t, reflect.DeepEqual(tt.wantErr.status, errStatus), "grpc status")
			} else if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			tt.want(got)
		})
	}
}

func setupForUpdate(t *testing.T) *testDeps {
	setup := newSetup(t, tfx)
	setup.HaveInDatabase("webhook", assertdb.Row{
		"id":         32767,
		"account_id": "42",
		"event":      "event1",
		"name":       "name1",
		"url":        "url1",
		"rate_limit": 2,
		"created_at": "2020-01-12 22:41:42",
		"updated_at": "2020-01-12 22:41:42"})

	return setup
}

func Test_Update(t *testing.T) {
	log.Println("test Update")

	i := setupForUpdate(t)
	defer i.teardown(t)
	client := i.getClient(t)

	type wantErr struct {
		status *status.Status
		ok     bool
	}
	tests := []struct {
		name    string
		params  *webhookpb.UpdateParams
		want    func(*webhookpb.UpdateReply)
		wantErr wantErr
	}{
		{
			name: "happy",
			params: &webhookpb.UpdateParams{
				Id:        int64(32767),
				AccountId: "42",
				Event:     "event2",
				Name:      "name2",
				URL:       "url2",
				RateLimit: int32(50),
			},
			want: func(response *webhookpb.UpdateReply) {

				assert.Equal(t, response.Webhook.GetId(), int64(32767), "Id")
				assert.Equal(t, response.Webhook.GetAccountId(), "42", "AccountId")
				assert.Equal(t, response.Webhook.GetEvent(), "event2", "Event")
				assert.Equal(t, response.Webhook.GetName(), "name2", "Name")
				assert.Equal(t, response.Webhook.GetURL(), "url2", "URL")
				assert.Equal(t, response.Webhook.GetRateLimit(), int32(50), "RateLimit")
				assert.Equal(t, response.Webhook.GetCreatedAt().AsTime().Format(assertdb.SQLTimestampWithoutTimeZone), "2020-01-12 22:41:42", "CreatedAt")
				assert.Check(t, response.Webhook.GetUpdatedAt().AsTime().Format(assertdb.SQLTimestampWithoutTimeZone) > "2020-01-12 22:41:42", "UpdatedAt")
				i.SeeInDatabase("webhook", assertdb.Criteria{
					"id":           32767,
					"account_id":   "42",
					"event":        "event2",
					"name":         "name2",
					"url":          "url2",
					"created_at":   "2020-01-12 22:41:42",
					"updated_at >": "2020-01-12 22:41:42",
				})

			},
		},
		{
			name: "not found id",
			params: &webhookpb.UpdateParams{
				Id:        int64(32776),
				AccountId: "42",
				Event:     "event2",
				Name:      "name2",
				URL:       "url2",
				RateLimit: int32(50),
			},
			want: func(response *webhookpb.UpdateReply) {
				i.DontSeeInDatabase("webhook", assertdb.Criteria{
					"id": 32776,
				})
			},
			wantErr: wantErr{
				status: status.New(codes.Unknown, "not found"),
				ok:     true, // the grpc service did respond
			},
		},
		{
			name: "not found account id",
			params: &webhookpb.UpdateParams{
				Id:        int64(32767),
				AccountId: "43",
				Event:     "event2",
				Name:      "name2",
				URL:       "url2",
				RateLimit: int32(50),
			},
			want: func(response *webhookpb.UpdateReply) {
				i.DontSeeInDatabase("webhook", assertdb.Criteria{
					"id":         32767,
					"account_id": "43",
				})
			},
			wantErr: wantErr{
				status: status.New(codes.Unknown, "not found"),
				ok:     true, // the grpc service did respond
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.Update(i.ctx, tt.params)
			if tt.wantErr.status != nil && err != nil {
				errStatus, ok := status.FromError(err)
				assert.Equal(t, ok, tt.wantErr.ok, "grpc ok")
				assert.Check(t, reflect.DeepEqual(tt.wantErr.status, errStatus), "grpc status")
			} else if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			tt.want(got)
		})
	}

}

func setupForDelete(t *testing.T) *testDeps {
	setup := newSetup(t, tfx)
	setup.HaveInDatabase("webhook", assertdb.Row{
		"id":         32767,
		"account_id": "42",
		"event":      "event1",
		"name":       "name1",
		"url":        "url1",
		"rate_limit": 2,
		"created_at": "2021-01-12 22:41:42",
		"updated_at": "2021-01-13 22:25:25"})
	return setup
}

func Test_Delete(t *testing.T) {
	log.Println("test Delete")

	i := setupForDelete(t)
	defer i.teardown(t)
	client := i.getClient(t)

	type wantErr struct {
		status *status.Status
		ok     bool
	}
	tests := []struct {
		name    string
		params  *webhookpb.DeleteParams
		want    func(*webhookpb.NoReply)
		wantErr wantErr
	}{
		{
			name: "happy",
			params: &webhookpb.DeleteParams{
				Id:        32767,
				AccountId: "42",
			},
			want: func(*webhookpb.NoReply) {
				i.DontSeeInDatabase("webhook", assertdb.Criteria{
					"id": 32767,
				})
			},
		},

		{
			name: "not found",
			params: &webhookpb.DeleteParams{
				Id:        32777,
				AccountId: "42",
			},
			want: func(*webhookpb.NoReply) {
				i.DontSeeInDatabase("webhook", assertdb.Criteria{
					"id": 32777,
				})
			},
			wantErr: wantErr{
				status: status.New(codes.Unknown, "not found"),
				ok:     true, // the grpc service did respond
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.Delete(i.ctx, tt.params)
			if tt.wantErr.status != nil && err != nil {
				errStatus, ok := status.FromError(err)
				assert.Equal(t, ok, tt.wantErr.ok, "grpc ok")
				assert.Check(t, reflect.DeepEqual(tt.wantErr.status, errStatus), "grpc status")
			} else if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			tt.want(got)
		})
	}
}
