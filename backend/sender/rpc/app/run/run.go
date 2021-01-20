package run

import (
	"github.com/burstsms/mtmo-tp/backend/lib/rpcbuilder"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/app/db"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/app/service"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
)

func Server(deps rpcbuilder.Deps) error {
	pdb := db.NewSQLDB(deps.PostgresConn)
	senderpb.RegisterServiceServer(deps.Server, service.NewSenderService(pdb))
	return nil
}
