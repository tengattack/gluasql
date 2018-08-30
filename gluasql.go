package gluasql

import (
	mysql "github.com/tengattack/gluasql/mysql"
	"github.com/yuin/gopher-lua"
)

func Preload(L *lua.LState) {
	L.PreloadModule("mysql", mysql.Loader)
}
