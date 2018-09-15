package gluasql_sqlite3

import (
	"database/sql"
	"time"

	lua "github.com/yuin/gopher-lua"
)

const (
	CLIENT_TYPENAME = "sqlite3{client}"
)

// Client sqlite3
type Client struct {
	DB      *sql.DB
	Timeout time.Duration
}

var clientMethods = map[string]lua.LGFunction{
	"open":        clientOpenMethod,
	"set_timeout": clientSetTimeoutMethod,
	"close":       clientCloseMethod,
	"query":       clientQueryMethod,
}

func checkClient(L *lua.LState) *Client {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Client); ok {
		return v
	}
	L.ArgError(1, "client expected")
	return nil
}

func clientCloseMethod(L *lua.LState) int {
	client := checkClient(L)

	if client.DB == nil {
		L.Push(lua.LBool(true))
		return 1
	}

	err := client.DB.Close()
	// always clean
	client.DB = nil
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}
