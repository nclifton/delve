package db

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/burstsms/mtmo-tp/backend/lib/assertdb"
	"github.com/burstsms/mtmo-tp/backend/lib/fixtures"
)

func Test_sqlDB_SenderAddressExists(t *testing.T) {

	// test fixture - require a postgres instance with the sender database, and the sender table containing a sender
	tfx := fixtures.New(fixtures.Config{
		Name:          "sender_sql",
		MigrationRoot: "file://../../../sender/migration/sql",
	})
	tfx.SetupPostgres("sender")
	s := newSetup(t, tfx)

	s.HaveInDatabase("sender", assertdb.Row{
		"id":               s.uuids[0],
		"account_id":       s.uuids[1],
		"address":          "FISH",
		"mms_provider_key": "optus",
		"channels":         []string{"mms", "sms"},
		"country":          "AU",
		"comment":          "Slartibartfast",
		"created_at":       s.dates[0],
		"updated_at":       s.dates[0]})

	postgresConn, err := pgxpool.Connect(s.ctx, tfx.Postgres.ConnStr)
	if err != nil {
		t.Fatalf("failed to get postgres connection pool for test: %s", err.Error())
	}
	closedPostgresConn, err := pgxpool.Connect(s.ctx, tfx.Postgres.ConnStr)
	if err != nil {
		t.Fatalf("failed to get postgres connection pool for test: %s", err.Error())
	}
	closedPostgresConn.Close()

	type fields struct {
		db *pgxpool.Pool
	}

	type args struct {
		ctx     context.Context
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "sender does exist",
			fields:  fields{postgresConn},
			args:    args{s.ctx, "FISH"},
			want:    true,
			wantErr: false,
		}, {
			name:    "sender does not exist",
			fields:  fields{postgresConn},
			args:    args{s.ctx, "CHIPS"},
			want:    false,
			wantErr: false,
		}, {
			name:    "sender does not exist",
			fields:  fields{closedPostgresConn},
			args:    args{s.ctx, "CHIPS"},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &sqlDB{
				tt.fields.db,
			}
			got, err := db.SenderAddressExists(tt.args.ctx, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("sqlDB.SenderAddressExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("sqlDB.SenderAddressExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testDeps struct {
	assertdb.AssertDb
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
	date, _ := time.Parse(assertdb.SQLTimestampWithoutTimeZone, "2020-01-13 14:52:21")
	for i := range dates {
		dates[i] = date.Add(time.Duration(-24*i) * time.Hour)
	}

	return &testDeps{
		ctx:      context.Background(),
		AssertDb: *assertdb.New(t, tfx.Postgres.ConnStr),
		uuids:    uuids,
		dates:    dates,
	}

}
