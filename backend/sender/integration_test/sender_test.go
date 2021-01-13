// +build integration

package test

import (
	"os"
	"testing"

	"github.com/burstsms/mtmo-tp/backend/lib/fixtures"
	"github.com/burstsms/mtmo-tp/backend/sender/rpc/app/run"
)

var tfx *fixtures.TestFixtures

func TestMain(m *testing.M) {
	tfx = fixtures.New()
	tfx.SetupPostgres("sender")
	tfx.GRPCStart(run.Server)
	code := m.Run()
	defer os.Exit(code)
	defer tfx.Teardown()
}
