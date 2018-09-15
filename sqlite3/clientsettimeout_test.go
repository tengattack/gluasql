package gluasql_sqlite3_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tengattack/gluasql"
	"github.com/yuin/gopher-lua"
)

func TestClientSetTimeout(t *testing.T) {
	assert := assert.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluasql.Preload(L)

	script := `
		c=require 'sqlite3'.new();
		c:set_timeout(1000);
	`

	assert.NoError(L.DoString(script))
}
