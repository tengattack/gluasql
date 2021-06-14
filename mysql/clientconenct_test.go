package gluasql_mysql_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"

	"github.com/tengattack/gluasql"
	gluamysql "github.com/tengattack/gluasql/mysql"
	util "github.com/tengattack/gluasql/util"
)

func TestParseConnectionString(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	timeout, err := time.ParseDuration("2s")
	require.NoError(err)

	scheme, dsn, err := gluamysql.ParseConnectionString(nil, timeout)
	require.Error(err)
	assert.Equal(gluamysql.ErrConnectionString, err)
	assert.Empty(scheme)
	assert.Empty(dsn)

	connectionString := "user:123456@tcp(127.0.0.1:3306)/foo"
	scheme, dsn, err = gluamysql.ParseConnectionString(connectionString, timeout)
	require.NoError(err)
	assert.Equal("mysql", scheme)
	assert.Equal(connectionString+"?readTimeout=2s&writeTimeout=2s", dsn)

	connectionString = "user:123456@tcp(127.0.0.1:3306)/foo?charset=utf8mb4"
	scheme, dsn, err = gluamysql.ParseConnectionString(connectionString, timeout)
	require.NoError(err)
	assert.Equal("mysql", scheme)
	assert.Equal(connectionString+"&readTimeout=2s&writeTimeout=2s", dsn)

	cs := map[string]interface{}{
		"user":     "user",
		"password": "123456",
		"database": "foo",
		"charset":  "utf8mb4",
	}
	scheme, dsn, err = gluamysql.ParseConnectionString(cs, timeout)
	require.NoError(err)
	assert.Equal("mysql", scheme)
	assert.Equal(connectionString+"&readTimeout=2s&writeTimeout=2s", dsn)

	connectionString = "postgres://host=127.0.0.1 port=5432 user=user " +
		"password=123456 dbname=foo sslmode=disable"
	scheme, dsn, err = gluamysql.ParseConnectionString(connectionString, timeout)
	require.NoError(err)
	assert.Equal("postgres", scheme)
	assert.Equal(connectionString[11:]+" connect_timeout=2", dsn)
}

func TestClientConnect(t *testing.T) {
	assert := assert.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluasql.Preload(L)

	script := getLuaDbConnection() + `
		return ok, err;
	`
	assert.NoError(L.DoString(script))

	ok := util.GetValue(L, 1)
	err := util.GetValue(L, 2)
	assert.True(ok.(bool))
	assert.Nil(err)
}
