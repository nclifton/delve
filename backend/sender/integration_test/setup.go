// +build integration

package test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/burstsms/mtmo-tp/backend/lib/assertdb"
	"github.com/burstsms/mtmo-tp/backend/lib/fixtures"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/senderpb"
)

type testDeps struct {
	assertdb.AssertDb // could this be moved to fixtures?
	fixtures.TestFixtures
	ctx   context.Context
	uuids []string
	dates []time.Time
}

func newSetup(t *testing.T, tfx *fixtures.TestFixtures) *testDeps {

	uuids := make([]string, 20)
	for i := range uuids {
		uuids[i] = uuid.New().String()
	}
	dates := make([]time.Time, 20)
	date , _:= time.Parse(assertdb.SQLTimestampWithoutTimeZone, "2020-01-13 14:52:21")
	for i := range dates {
		dates[i] = date.Add(time.Duration(-24 * i) * time.Hour)
	}

	return &testDeps{
		ctx:          context.Background(),
		TestFixtures: *tfx,
		AssertDb:     *assertdb.New(t, tfx.Postgres.ConnStr),
		uuids:        uuids,
		dates:        dates,
	}

}

func (s *testDeps) teardown(t *testing.T) {
	s.AssertDb.Teardown()
}

func (s *testDeps) getClient(t *testing.T) senderpb.ServiceClient {
	return senderpb.NewServiceClient(s.GRPCClientConnection(t, s.ctx))
}
