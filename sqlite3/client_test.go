package gluasql_sqlite3_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/tengattack/gluasql"
	lua "github.com/yuin/gopher-lua"
)

var (
	timeNow time.Time
)

func TestMain(m *testing.M) {
	// prepare test db
	err := createTestDatabase()
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
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
		c=require 'sqlite3'.new();
		c:close();
	`)
	assert.NoError(L.DoString(script))
}

func createTestDatabase() error {
	const (
		tableName = "mytable"
	)

	db, err := sql.Open("sqlite3", "file:test.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// clean up
	_, err = db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s`, tableName))
	if err != nil {
		return err
	}

	_, err = db.Exec(fmt.Sprintf(`CREATE TABLE %s (
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			phone_numbers TEXT NOT NULL,
			created_at TEXT NOT NULL
		)`, tableName))
	if err != nil {
		return err
	}

	timeNow = time.Now()
	timeNowFormat := timeNow.UTC().Format("2006-01-02 15:04:05")
	rows := [][]string{
		[]string{"John Doe", "john@doe.com", `["555-555-555"]`, timeNowFormat},
		[]string{"John Doe", "johnalt@doe.com", `[]`, timeNowFormat},
		[]string{"Jane Doe", "jane@doe.com", `[]`, timeNowFormat},
		[]string{"Evil Bob", "evilbob@gmail.com", `["555-666-555","666-666-666"]`, timeNowFormat},
	}

	args := make([]interface{}, len(rows)*len(rows[0]))
	sql := fmt.Sprintf("INSERT INTO %s (name, email, phone_numbers, created_at) VALUES\n", tableName)
	i := 0
	for j, row := range rows {
		sql += "(?, ?, ?, ?)"
		if j+1 < len(rows) {
			sql += ",\n"
		} else {
			sql += ";"
		}
		for k, val := range row {
			args[i+k] = val
		}
		i += len(row)
	}
	_, err = db.Exec(sql, args...)
	return err
}

func getLuaDbConnection() string {
	return fmt.Sprintf(`
		c=require 'sqlite3'.new();
		ok, err = c:open("%s", { cache = "%s", mode = "%s" });
	`, "test.db", "shared", "ro")
}
