// +build integration

package test

import (
	"context"
	"testing"
	"time"

	"github.com/burstsms/mtmo-tp/backend/lib/assertdb"
	"github.com/burstsms/mtmo-tp/backend/lib/asserthttp"
	"github.com/burstsms/mtmo-tp/backend/lib/fixtures"
	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
)

var testStartTime time.Time

func init() {
	testStartTime = time.Now()
}

type testDeps struct {
	assertdb.AssertDb // could this be moved to fixtures
	fixtures.TestFixtures
	// why can't I embed the http receiver as well?
	// could http receiver be supported in the fixtures?
	httpReceiver      *asserthttp.HttpReceiver
	ctx               context.Context
	webhookURL        string
	GetURL            func() string
	ResetHttpRequests func()
	WaitForRequests   func(ExpectedRequests)
	AssertRequests    func(ExpectedRequests)
}

type ExpectedRequests = asserthttp.ExpectedRequests

func jsonString(t *testing.T, v interface{}) string {
	return asserthttp.JSONString(t, v)
}

func newSetup(t *testing.T, tfx *fixtures.TestFixtures) *testDeps {
	s := &testDeps{
		TestFixtures: *tfx,
		AssertDb:     *assertdb.New(t, tfx.Postgres.ConnStr),
		httpReceiver: asserthttp.NewHttpReceiver(t),
		ctx:          context.Background(),
	}
	s.GetURL = s.httpReceiver.GetURL
	s.webhookURL = s.GetURL()
	s.ResetHttpRequests = s.httpReceiver.ResetHttpRequests
	s.WaitForRequests = s.httpReceiver.WaitForRequests
	s.AssertRequests = s.httpReceiver.AssertRequests
	return s
}

func (setup *testDeps) getClient(t *testing.T) webhookpb.ServiceClient {
	return webhookpb.NewServiceClient(setup.GRPCClientConnection(t, setup.ctx))
}

func (setup *testDeps) teardown(t *testing.T) {
	setup.AssertDb.Teardown()
	setup.httpReceiver.Teardown()

}
