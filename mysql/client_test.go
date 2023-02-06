package gluasql_mysql_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"
	"github.com/stretchr/testify/assert"
	"github.com/tengattack/gluasql"

	lua "github.com/yuin/gopher-lua"
)

var (
	timeNow   time.Time
	dbName    = "mydb"
	tableName = "mytable"
	address   = "localhost"
	port      = 3306
)

func TestMain(m *testing.M) {
	// prepare mysql server
	ctx := sql.NewEmptyContext()
	engine := sqle.NewDefault(
		memory.NewDBProvider(
			createTestDatabase(ctx),
		))

	config := server.Config{
		Protocol: "tcp",
		Address:  fmt.Sprintf("%s:%d", address, port),
	}

	s, err := server.NewDefaultServer(config, engine)
	if err != nil {
		panic(err)
	}

	go s.Start()
	code := m.Run()
	s.Close()

	os.Exit(code)
}

func TestClientSetKeepalive(t *testing.T) {
	assert := assert.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluasql.Preload(L)

	script := getLuaDbConnection() + `
		c:set_keepalive(2000, 5);
	`
	assert.NoError(L.DoString(script))
}
func TestClientClose(t *testing.T) {
	assert := assert.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluasql.Preload(L)

	script := getLuaDbConnection() + `
		c:close();
	`
	assert.NoError(L.DoString(script))

	script = fmt.Sprintf(`
		c=require 'mysql'.new();
		c:close();
	`)
	assert.NoError(L.DoString(script))
}

func createTestDatabase(ctx *sql.Context) *memory.Database {
	db := memory.NewDatabase(dbName)
	table := memory.NewTable(tableName, sql.NewPrimaryKeySchema(sql.Schema{
		{Name: "name", Type: types.Text, Nullable: false, Source: tableName},
		{Name: "email", Type: types.Text, Nullable: false, Source: tableName},
		{Name: "phone_numbers", Type: types.JSON, Nullable: false, Source: tableName},
		{Name: "created_at", Type: types.Timestamp, Nullable: false, Source: tableName},
	}), db.GetForeignKeyCollection())
	db.AddTable(tableName, table)

	creationTime := time.Unix(0, 1667304000000001000).UTC()
	_ = table.Insert(ctx, sql.NewRow("Jane Deo", "janedeo@gmail.com", types.MustJSON(`["556-565-566", "777-777-777"]`), creationTime))
	_ = table.Insert(ctx, sql.NewRow("Jane Doe", "jane@doe.com", types.MustJSON(`[]`), creationTime))
	_ = table.Insert(ctx, sql.NewRow("John Doe", "john@doe.com", types.MustJSON(`["555-555-555"]`), creationTime))
	_ = table.Insert(ctx, sql.NewRow("John Doe", "johnalt@doe.com", types.MustJSON(`[]`), creationTime))
	return db
}

func getLuaDbConnection() string {
	return fmt.Sprintf(`
		c=require 'mysql'.new();
		ok, err = c:connect({ host = "%s", port = %d, database = "%s", user = "%s", password = "%s" });
	`, address, port, dbName, "user", "pass")
}
