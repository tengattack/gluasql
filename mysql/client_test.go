package gluasql_mysql_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tengattack/gluasql"
	"github.com/yuin/gopher-lua"

	"gopkg.in/src-d/go-mysql-server.v0"
	"gopkg.in/src-d/go-mysql-server.v0/mem"
	"gopkg.in/src-d/go-mysql-server.v0/server"
	"gopkg.in/src-d/go-mysql-server.v0/sql"
	"gopkg.in/src-d/go-vitess.v0/mysql"
)

var (
	timeNow time.Time
)

func TestMain(m *testing.M) {
	// prepare mysql server
	driver := sqle.NewDefault()
	driver.AddDatabase(createTestDatabase())

	auth := mysql.NewAuthServerStatic()
	auth.Entries["user"] = []*mysql.AuthServerStaticEntry{{
		Password: "pass",
	}}

	config := server.Config{
		Protocol: "tcp",
		Address:  "localhost:3306",
		Auth:     auth,
	}

	s, err := server.NewDefaultServer(config, driver)
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

func createTestDatabase() *mem.Database {
	const (
		dbName    = "test"
		tableName = "mytable"
	)

	db := mem.NewDatabase(dbName)
	table := mem.NewTable(tableName, sql.Schema{
		{Name: "name", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "email", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "phone_numbers", Type: sql.JSON, Nullable: false, Source: tableName},
		{Name: "created_at", Type: sql.Timestamp, Nullable: false, Source: tableName},
	})

	db.AddTable(tableName, table)
	ctx := sql.NewEmptyContext()

	timeNow = time.Now()
	rows := []sql.Row{
		sql.NewRow("John Doe", "john@doe.com", []string{"555-555-555"}, timeNow),
		sql.NewRow("John Doe", "johnalt@doe.com", []string{}, timeNow),
		sql.NewRow("Jane Doe", "jane@doe.com", []string{}, timeNow),
		sql.NewRow("Evil Bob", "evilbob@gmail.com", []string{"555-666-555", "666-666-666"}, timeNow),
	}

	for _, row := range rows {
		table.Insert(ctx, row)
	}

	return db
}

func getLuaDbConnection() string {
	return fmt.Sprintf(`
		c=require 'mysql'.new();
		ok, err = c:connect({ host = "%s", port = %d, database = "%s", user = "%s", password = "%s" });
	`, "localhost", 3306, "test", "user", "pass")
}
