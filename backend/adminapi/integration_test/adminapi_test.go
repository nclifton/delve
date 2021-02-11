// +build integration

package test

import (
	"os"
	"testing"

	"github.com/burstsms/mtmo-tp/backend/lib/fixtures"
)

var tfx *fixtures.TestFixtures

func TestMain(m *testing.M) {
	setupFixtures()
	code := m.Run()
	defer os.Exit(code)
	defer tfx.Teardown()
}

func setupFixtures(){
	tfx = fixtures.New(fixtures.Config{Name: "adminapi"})
	tfx.SetupPostgres("sender")
}