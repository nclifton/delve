package assertdb

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgx/v4"
	"gotest.tools/assert"
)

const SQLDateTime = "2006-01-02 15:04:05"

type AssertDb struct {
	t               *testing.T
	connStr         string
	haveTableRowIds map[string][]interface{}
	ctx             context.Context
	conn            *pgx.Conn
}

func New(t *testing.T, connStr string) *AssertDb {
	adb := &AssertDb{
		t:               t,
		connStr:         connStr,
		haveTableRowIds: make(map[string][]interface{}),
	}
	adb.getConnection()
	return adb
}

type criterion []interface{}

type Criteria []criterion

type value struct {
	Name  string
	Value string
}

type Values []value

func (adb *AssertDb) SeeInDatabase(table string, criteria Criteria) {
	count := adb.countInDatabase(table, criteria)
	assert.Check(adb.t, count > 0, fmt.Sprintf("see in table %s where %+v", table, criteria))
}

func (adb *AssertDb) DontSeeInDatabase(table string, criteria Criteria) {
	count := adb.countInDatabase(table, criteria)
	assert.Check(adb.t, count == 0, fmt.Sprintf("see in table %s where %+v", table, criteria))
}

func (adb *AssertDb) countInDatabase(table string, criteria Criteria) int {

	where := make([]string, len(criteria))
	values := make([]interface{}, len(criteria))
	for idx, c := range criteria {
		name := fmt.Sprintf("%v", c[0])
		op := "="
		if strings.Contains("=<>", string(name[len(name)-1:])) {
			op = ""
		}
		where[idx] = fmt.Sprintf("%s %s $%d", name, op, idx+1)
		values[idx] = c[1]
	}

	row := adb.conn.QueryRow(adb.ctx,
		fmt.Sprintf(`SELECT COUNT(*) FROM "%s" WHERE %s`, table, strings.Join(where, " AND ")), values...)

	var count int
	assert.NilError(adb.t, row.Scan(&count), fmt.Sprintf("query table %s failed", table))

	return count
}

func (adb *AssertDb) getConnection() {
	adb.ctx = context.Background()
	conn, err := pgx.Connect(adb.ctx, adb.connStr)
	if err != nil {
		adb.t.Fatalf("Unable to connect to %v error: %+v\n", adb.connStr, err)
	}
	adb.conn = conn
}

func (adb *AssertDb) HaveInDatabase(table string, names string, arguments []interface{}) {

	values := make([]string, len(arguments))
	for idx := range arguments {
		values[idx] = fmt.Sprintf("$%d", idx+1)
	}

	sql := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s) RETURNING id`,
		table, names, strings.Join(values, ","))

	var insertId interface{}
	err := adb.conn.QueryRow(adb.ctx, sql, arguments...).Scan(&insertId)
	if err != nil {
		adb.t.Fatalf("Unable to insert into table %s\nerror: %+v\n", table, err)
	}
	adb.haveTableRowIds[table] = append(adb.haveTableRowIds[table], insertId)
}

// after calling this, a new AssertDb will need to be created
func (adb *AssertDb) Teardown() {
	for table, ids := range adb.haveTableRowIds {
		for _, id := range ids {
			sql := fmt.Sprintf(`DELETE FROM "%s" WHERE id = $1`, table)
			_, err := adb.conn.Exec(adb.ctx, sql, id)
			if err != nil {
				adb.t.Fatalf("Unable to delete have in database items from table %s\nerror: %+v\n", table, err)
			}
		}
	}
	if !adb.conn.IsClosed() {
		adb.conn.Close(adb.ctx)
	}
	adb = nil
}
func (adb *AssertDb) Close() {
	adb.Teardown()
}
