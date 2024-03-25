package gluasql_mysql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tengattack/gluasql"
	util "github.com/tengattack/gluasql/util"
	"github.com/yuin/gopher-lua"
)

var expectedResults = []interface{}{
	map[string]interface{}{
		"created_at":    "2022-11-01 12:00:00.000001",
		"email":         "janedeo@gmail.com",
		"name":          "Jane Deo",
		"phone_numbers": "[\"556-565-566\",\"777-777-777\"]"},
	map[string]interface{}{
		"created_at":    "2022-11-01 12:00:00.000001",
		"email":         "jane@doe.com",
		"name":          "Jane Doe",
		"phone_numbers": "[]"},
}

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
	// sql.NewRow("Jane Deo", "janedeo@gmail.com", types.MustJSON(`["556-565-566", "777-777-777"]`), creationTime)
	// sql.NewRow("Jane Doe", "jane@doe.com", types.MustJSON(`[]`), creationTime)
	assert.Equal(expectedResults, res)
}
