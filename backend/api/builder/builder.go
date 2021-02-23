package builder

import (
	"context"
	"log"

	"github.com/burstsms/mtmo-tp/backend/account/rpc/accountpb"
	"github.com/burstsms/mtmo-tp/backend/api/service"
	"github.com/burstsms/mtmo-tp/backend/lib/rest"
	"github.com/burstsms/mtmo-tp/backend/lib/restbuilder"
	"github.com/burstsms/mtmo-tp/backend/lib/rpcbuilder"
	"github.com/burstsms/mtmo-tp/backend/lib/valid"
	mms "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
	sms "github.com/burstsms/mtmo-tp/backend/sms/rpc/client"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
	"github.com/kelseyhightower/envconfig"
)

func NewFromEnv() *serviceBuilder {
	var conf Config
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Fatal("Failed to read env vars:", err)
	}
	return &serviceBuilder{conf: conf}
}

type Config struct {
	AccountRPCAddress string `envconfig:"ACCOUNT_RPC_ADDRESS"`
	SenderRPCAddress  string `envconfig:"SENDER_RPC_ADDRESS"`
	WebhookRPCAddress string `envconfig:"WEBHOOK_RPC_ADDRESS"`
	MMSRPCAddress     string `envconfig:"MMS_RPC_ADDRESS"`
	SMSRPCAddress     string `envconfig:"SMS_RPC_ADDRESS"`
}

type serviceBuilder struct {
	clients       *service.Clients
	authenticator rest.Authenticator
	conf          Config
}

func (b *serviceBuilder) SetClients(clients *service.Clients) {
	b.clients = clients
}

func (b *serviceBuilder) SetAuthenticator(authenticator rest.Authenticator) {
	b.authenticator = authenticator
}

func (b *serviceBuilder) Run(deps restbuilder.Deps) error {
	ctx := context.Background()

	if b.clients == nil {
		b.clients = &service.Clients{
			SenderClient: senderpb.NewServiceClient(
				rpcbuilder.NewClientConn(b.conf.SenderRPCAddress, deps.Tracer),
			),
			WebhookClient: webhookpb.NewServiceClient(
				rpcbuilder.NewClientConn(b.conf.WebhookRPCAddress, deps.Tracer),
			),
			AccountClient: accountpb.NewServiceClient(
				rpcbuilder.NewClientConn(b.conf.AccountRPCAddress, deps.Tracer),
			),
			MMSClient: mms.New(b.conf.MMSRPCAddress),
			SMSClient: sms.New(b.conf.SMSRPCAddress),
		}
	}

	if b.authenticator == nil {
		b.authenticator = func(key string) (interface{}, error) {
			reply, err := b.clients.AccountClient.FindAccountByAPIKey(ctx, &accountpb.FindAccountByAPIKeyParams{Key: key})
			if err != nil {
				return nil, err
			}

			return reply.Account, nil
		}
	}

	hb := rest.NewHandlerBuilder(&rest.HandlerConfig{
		Log:           deps.Log,
		Tracer:        deps.Tracer,
		JSONValidator: valid.Validate,
	})

	authRoute := hb().SetMiddleware(
		rest.NewAuthMiddleware(b.authenticator),
		rest.NewTracingMiddleware(),
		rest.NewLoggingMiddleware(),
	).Handle
	baseRoute := hb().SetMiddleware(
		rest.NewTracingMiddleware(),
		rest.NewLoggingMiddleware(),
	).Handle

	service.Routes(deps.Router, baseRoute, authRoute, b.clients)

	return nil
}
