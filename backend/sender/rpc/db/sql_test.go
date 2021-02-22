package db

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/burstsms/mtmo-tp/backend/lib/assertdb"
	"github.com/burstsms/mtmo-tp/backend/lib/fixtures"
)

var tfx *fixtures.TestFixtures

func TestMain(m *testing.M) {
	// test fixture - require a postgres instance with the sender database
	tfx = fixtures.New(fixtures.Config{
		Name:          "sender_sql",
		MigrationRoot: "file://../../../sender/migration/sql",
	})
	tfx.SetupPostgres("sender")
	code := m.Run()
	defer os.Exit(code)
	defer tfx.Teardown()
}

func Test_sqlDB_SenderAddressExists(t *testing.T) {

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
			fields:  fields{s.postgresConn},
			args:    args{s.ctx, "FISH"},
			want:    true,
			wantErr: false,
		}, {
			name:    "sender does not exist",
			fields:  fields{s.postgresConn},
			args:    args{s.ctx, "CHIPS"},
			want:    false,
			wantErr: false,
		}, {
			name:    "sender does not exist",
			fields:  fields{s.closedPostgresConn},
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

func Test_sqlDB_GetSenderEnums(t *testing.T) {

	s := newSetup(t, tfx)

	type fields struct {
		sql *pgxpool.Pool
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    SenderEnums
		wantErr bool
	}{
		{
			name:   "happy",
			fields: fields{s.postgresConn},
			args:   args{s.ctx},
			want: SenderEnums{
				"provider_key": []string{"fake", "optus", "mgage"},
				"channel":      []string{"mms", "sms"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &sqlDB{
				sql: tt.fields.sql,
			}
			got, err := db.GetSenderEnums(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("sqlDB.GetSenderEnums() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sqlDB.GetSenderEnums() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testDeps struct {
	*assertdb.AssertDb
	*fixtures.TestFixtures
	ctx                context.Context
	uuids              []string
	dates              []time.Time
	postgresConn       *pgxpool.Pool
	closedPostgresConn *pgxpool.Pool
}

func newSetup(t *testing.T, tfx *fixtures.TestFixtures) *testDeps {

	s := &testDeps{
		AssertDb:     assertdb.New(t, tfx.Postgres.ConnStr),
		TestFixtures: tfx,
		ctx:          context.Background(),
		uuids:        make([]string, 20),
		dates:        make([]time.Time, 20),
	}

	for i := range s.uuids {
		s.uuids[i] = uuid.New().String()
	}
	date, _ := time.Parse(assertdb.SQLTimestampWithoutTimeZone, "2020-01-13 14:52:21")
	for i := range s.dates {
		s.dates[i] = date.Add(time.Duration(-24*i) * time.Hour)
	}

	var err error
	s.postgresConn, err = pgxpool.Connect(s.ctx, tfx.Postgres.ConnStr)
	if err != nil {
		t.Fatalf("failed to get postgres connection pool for test: %s", err.Error())
	}
	s.closedPostgresConn, err = pgxpool.Connect(s.ctx, tfx.Postgres.ConnStr)
	if err != nil {
		t.Fatalf("failed to get postgres connection pool for test: %s", err.Error())
	}
	s.closedPostgresConn.Close()

	return s

}
