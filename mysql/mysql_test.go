package gluasql_mysql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tengattack/gluasql"
	util "github.com/tengattack/gluasql/util"
	"github.com/yuin/gopher-lua"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluasql.Preload(L)

	script := `
		mysql=require 'mysql';
		return mysql.new();
	`

	assert.NoError(L.DoString(script))

	c := L.Get(1)
	assert.NotNil(c)
}

func TestEscape(t *testing.T) {
	assert := assert.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluasql.Preload(L)

	script := `
		mysql=require 'mysql';
		return mysql.escape('I\'m busy.\r\n"Apple"');
	`

	assert.NoError(L.DoString(script))

	s := util.GetValue(L, 1)
	assert.Equal(`I\'m busy.\r\n\"Apple\"`, s)
}
