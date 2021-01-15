/**
 * "tools" that can be used for "domain" or integration testing where a database needs to be setup and used for assertions
 */

package assertdb

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgx/v4"
	"gotest.tools/assert"
)

const SQLTimestampWithoutTimeZone = "2006-01-02 15:04:05"

type Row = map[string]interface{}
type Criteria = map[string]interface{}

type AssertDb struct {
	t               *testing.T
	connStr         string
	haveTableRowIds map[string][]interface{}
	ctx             context.Context
	conn            *pgx.Conn
}

/**
 * Establishes a new database connection using the provided connection string.
 * Is intended for use within a go test, hence the pointer to the current test is required for perfroming assertions
 */
func New(t *testing.T, connStr string) *AssertDb {
	adb := &AssertDb{
		t:               t,
		connStr:         connStr,
		haveTableRowIds: make(map[string][]interface{}),
	}
	adb.getConnection()
	return adb
}

type value struct {
	Name  string
	Value string
}

type Values []value

/**
 * "assert" that the database table contains rows matching the provided criteria.
 * Criteria is a map. The "key" is the field name and can be suffixed with a comparison operator.
 * Operator must be separated from the field name by a space.
 */
func (adb *AssertDb) SeeInDatabase(table string, criteria Criteria) {
	count := adb.CountInDatabase(table, criteria)
	assert.Check(adb.t, count > 0, fmt.Sprintf("see in table %s where %+v", table, criteria))
}

/**
 * "assert" that the database table does not contain row matching the provided criteria
 * Criteria is a map. The "key" is the field name and can be suffixed with a comparison operator.
 * Operator must be separated from the field name by a space.
 */
func (adb *AssertDb) DontSeeInDatabase(table string, criteria Criteria) {
	count := adb.CountInDatabase(table, criteria)
	assert.Check(adb.t, count == 0, fmt.Sprintf("see in table %s where %+v", table, criteria))
}

/**
 * Counts the number of records that match the provided criteria in the provided table name.
 */
func (adb *AssertDb) CountInDatabase(table string, criteria Criteria) int {

	where := make([]string, 0, len(criteria))
	values := make([]interface{}, 0, len(criteria))
	idx := 0
	for name, value := range criteria {
		idx++
		op := " ="
		if strings.Contains(strings.TrimSpace(string(name)), " ") {
			op = ""
		}
		where = append(where, fmt.Sprintf("%s%s $%d", name, op, idx))
		values = append(values, value)
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

/**
 * Inserts data into a database table during test setup.
 * Any rows inserted into the database by this function will be removed by the function Teardown.
 * Row is a map. Keys are the field names.
 */
func (adb *AssertDb) HaveInDatabase(table string, row map[string]interface{}) {

	prepValues := make([]string, 0, len(row))
	values := make([]interface{}, 0, len(row))
	names := make([]string, 0, len(row))
	idx := 0
	for name, value := range row {
		idx++
		names = append(names, name)
		prepValues = append(prepValues, fmt.Sprintf("$%d", idx))
		values = append(values, value)
	}

	sql := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s) RETURNING id`,
		table, strings.Join(names, ","), strings.Join(prepValues, ","))

	var insertId interface{}
	err := adb.conn.QueryRow(adb.ctx, sql, values...).Scan(&insertId)
	if err != nil {
		adb.t.Fatalf("Unable to insert into table %s\nerror: %+v\n", table, err)
	}
	adb.haveTableRowIds[table] = append(adb.haveTableRowIds[table], insertId)
}

/**
 * Removes any rows added to the database using the HaveInDatabase function..
 * Closes the database connection.
 * Note: after calling this, a new AssertDb will need to be created so that a new database connection is created.
 */
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
