package gluasql_sqlite3

import (
	"github.com/yuin/gopher-lua"
)

var exports = map[string]lua.LGFunction{
	"new": newFn,
}

func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)

	L.SetField(mod, "_DEBUG", lua.LBool(false))
	L.SetField(mod, "_VERSION", lua.LString("0.0.0"))

	registerClientType(L)

	return 1
}

func registerClientType(L *lua.LState) {
	mt := L.NewTypeMetatable(CLIENT_TYPENAME)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), clientMethods))
}

func newFn(L *lua.LState) int {
	client := &Client{}
	ud := L.NewUserData()
	ud.Value = client
	L.SetMetatable(ud, L.GetTypeMetatable(CLIENT_TYPENAME))
	L.Push(ud)
	return 1
}
