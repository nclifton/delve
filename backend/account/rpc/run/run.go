package run

import (
	"github.com/burstsms/mtmo-tp/backend/lib/rpcbuilder"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/db"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/service"
)

func Server(deps rpcbuilder.Deps) error {
	pdb := db.NewSQLDB(deps.PostgresConn)
	senderpb.RegisterServiceServer(deps.Server, service.NewSenderService(pdb))
	return nil
}
