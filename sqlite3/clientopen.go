package gluasql_sqlite3

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/yuin/gopher-lua"

	_ "github.com/mattn/go-sqlite3"
	util "github.com/tengattack/gluasql/util"
)

func clientOpenMethod(L *lua.LState) int {
	client := checkClient(L)
	file := L.ToString(2)
	tb := util.GetValue(L, 3)
	options, ok := tb.(map[string]interface{})

	if file == "" {
		L.ArgError(2, "file path string excepted")
		return 0
	}
	if tb == nil || !ok {
		L.ArgError(3, "options excepted")
		return 0
	}

	_auth, _ := options["_auth"].(bool)
	supportStringArgs := []string{"_auth_user", "_auth_pass", "_auth_crypt",
		"_loc", "_mutex", "_txlock", "cache", "mode"}
	q := url.Values{}
	for _, arg := range supportStringArgs {
		val, ok := options[arg].(string)
		if ok {
			q.Set(arg, val)
		}
	}

	if client.Timeout > 0 {
		stimeout := client.Timeout.String()
		q.Set("_timeout", stimeout)
	}

	dsn := fmt.Sprintf("file:%s", file)
	s := q.Encode()
	if s != "" {
		if _auth {
			dsn += "?_auth&" + s
		} else {
			dsn += "?" + s
		}
	} else if _auth {
		// REVIEW: _auth in DSN
		dsn += "?_auth"
	}

	var err error
	client.DB, err = sql.Open("sqlite3", dsn)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}
