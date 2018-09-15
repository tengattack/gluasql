package gluasql_sqlite3_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tengattack/gluasql"
	util "github.com/tengattack/gluasql/util"
	"github.com/yuin/gopher-lua"
)

func TestClientQuery(t *testing.T) {
	assert := assert.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluasql.Preload(L)

	script := getLuaDbConnection() + `
		res, err = c:query("SELECT * FROM mytable LIMIT 2");
		return res, err;
	`

	assert.NoError(L.DoString(script))

	res := util.GetValue(L, 1)
	err := util.GetValue(L, 2)
	assert.Nil(err)

	// source:
	// sql.NewRow("John Doe", "john@doe.com", []string{"555-555-555"}, timeNow)
	// sql.NewRow("John Doe", "johnalt@doe.com", []string{}, timeNow)
	timeNowFormat := timeNow.UTC().Format("2006-01-02 15:04:05")
	assert.Equal([]interface{}{
		map[string]interface{}{
			"email":         "john@doe.com",
			"phone_numbers": "[\"555-555-555\"]",
			"name":          "John Doe",
			"created_at":    timeNowFormat,
		},
		map[string]interface{}{
			"email":         "johnalt@doe.com",
			"phone_numbers": "[]",
			"name":          "John Doe",
			"created_at":    timeNowFormat,
		},
	}, res)
}
