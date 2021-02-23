package test

import (
	"context"
	"os"
	"testing"

	"github.com/burstsms/mtmo-tp/backend/account/rpc/accountpb"
	"github.com/burstsms/mtmo-tp/backend/account/rpc/builder"
	"github.com/burstsms/mtmo-tp/backend/account/rpc/db"
	"github.com/burstsms/mtmo-tp/backend/lib/fixtures"
	"github.com/jackc/pgx/v4"
)

var tfx *fixtures.TestFixtures

func TestMain(m *testing.M) {
	tfx = fixtures.New(fixtures.Config{Name: "account"})
	tfx.SetupPostgres("account")
	tfx.SetupRedis()

	tfx.GRPCStart(builder.NewBuilderFromEnv())
	code := m.Run()

	defer os.Exit(code)
	defer tfx.Teardown()
}

// TODO: rename?
type testDeps struct {
	DBConn        *pgx.Conn
	accountClient accountpb.ServiceClient
}

func newSetup(t *testing.T, tfx *fixtures.TestFixtures) *testDeps {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, tfx.Postgres.ConnStr)
	if err != nil {
		t.Fatalf("Postgres: unable to connect to %s, error:  %s", tfx.Postgres.ConnStr, err)
	}

	return &testDeps{
		DBConn:        conn,
		accountClient: accountpb.NewServiceClient(tfx.GRPCClientConnection(t, ctx)),
	}
}

///////////////////
// Equal
///////////////////

func AccountPBEqual(expected, actual *accountpb.Account) bool {
	if expected.Id != actual.Id {
		return false
	}
	if expected.CreatedAt.String() != actual.CreatedAt.String() {
		return false
	}
	if expected.UpdatedAt.String() != actual.UpdatedAt.String() {
		return false
	}
	if expected.Name != actual.Name {
		return false
	}
	if expected.AlarisUsername != actual.AlarisUsername {
		return false
	}
	if expected.AlarisPassword != actual.AlarisPassword {
		return false
	}
	if expected.AlarisUrl != actual.AlarisUrl {
		return false
	}
	return true
}

///////////////////
// db & cleanup
///////////////////

func insertAccounts(t *testing.T, db *pgx.Conn, accounts []db.Account) {
	for _, account := range accounts {
		if _, err := db.Exec(context.Background(),
			"INSERT INTO account (id, created_at, updated_at, name, alaris_username, alaris_password, alaris_url) values ($1, $2, $3, $4, $5, $6, $7)",
			account.ID,
			account.CreatedAt,
			account.UpdatedAt,
			account.Name,
			account.AlarisUsername,
			account.AlarisPassword,
			account.AlarisURL,
		); err != nil {
			t.Fatal(err)
		}
	}
}

func insertAccountAPIKeys(t *testing.T, db *pgx.Conn, accountAPIKeys []db.AccountAPIKey) {
	for _, accountAPIKey := range accountAPIKeys {
		if _, err := db.Exec(context.Background(),
			"INSERT INTO account_api_keys (id, account_id, description, key) values ($1, $2, $3, $4)",
			accountAPIKey.ID,
			accountAPIKey.AccountID,
			accountAPIKey.Description,
			accountAPIKey.Key,
		); err != nil {
			t.Fatal(err)
		}
	}
}

func cleanup(t *testing.T, db *pgx.Conn) {
	if _, err := db.Exec(context.Background(), "TRUNCATE TABLE account, account_api_keys CASCADE"); err != nil {
		t.Fatal(err)
	}
}
