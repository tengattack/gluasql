package gluasql_mysql

import (
	"reflect"

	"github.com/junhsieh/goexamples/fieldbinding/fieldbinding"
	util "github.com/tengattack/gluasql/util"
	"github.com/yuin/gopher-lua"
)

func clientQueryMethod(L *lua.LState) int {
	client := checkClient(L)
	query := L.ToString(2)

	if client.DB == nil {
		L.Push(lua.LBool(true))
		L.Push(lua.LString("connect required"))
		return 2
	}

	if query == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("query excepted"))
		return 2
	}

	rows, err := client.DB.Query(query)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer rows.Close()

	fb := fieldbinding.NewFieldBinding()
	cols, err := rows.Columns()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	fb.PutFields(cols)

	tb := L.NewTable()
	for rows.Next() {
		if err := rows.Scan(fb.GetFieldPtrArr()...); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		tbRow := util.ToTableFromMap(L, reflect.ValueOf(fb.GetFieldArr()))
		tb.Append(tbRow)
	}

	L.Push(tb)
	return 1
}
