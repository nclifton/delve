package builder

import (
	"context"

	"github.com/kelseyhightower/envconfig"

	"github.com/burstsms/mtmo-tp/backend/account/rpc/accountpb"
	"github.com/burstsms/mtmo-tp/backend/account/rpc/db"
	"github.com/burstsms/mtmo-tp/backend/account/rpc/service"
	"github.com/burstsms/mtmo-tp/backend/lib/logger"
	"github.com/burstsms/mtmo-tp/backend/lib/rpcbuilder"
)

type Config struct {
}

type builder struct {
	conf Config
}

func NewBuilder(config Config) *builder {
	return &builder{conf: config}
}

func NewBuilderFromEnv() *builder {
	stLog := logger.NewLogger()
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		stLog.Fatalf(context.Background(), "envconfig.Process", "failed to read env vars: %s", err)
	}

	return NewBuilder(config)
}

func (b *builder) Run(deps rpcbuilder.Deps) error {
	pdb := db.NewSQLDB(deps.PostgresConn)

	if err := deps.Redis.EnableCache(); err != nil {
		return err
	}

	accountpb.RegisterServiceServer(deps.Server, service.NewAccountService(pdb, deps.Redis))

	return nil
}
