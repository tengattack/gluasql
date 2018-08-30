package gluasql_util_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	util "github.com/tengattack/gluasql/util"
	"github.com/yuin/gopher-lua"
)

func testLuaState() *lua.LState {
	l := lua.NewState()
	return l
}

func TestGetArbitrary(t *testing.T) {
	l := testLuaState()
	b := bytes.NewBufferString(`
		return {
			a = 1,
			b = 1.1,
			c = "foo",
			d = {
				e = "baz",
			},
			f = {"buz", 5},
			g = {},
		}
	`)
	fn, err := l.Load(b, "")
	require.Nil(t, err)
	l.Push(fn)
	l.Call(0, 1)
	i := util.GetValue(l, 1)
	assert.Equal(t, map[string]interface{}{
		"a": 1,
		"b": 1.1,
		"c": "foo",
		"d": map[string]interface{}{
			"e": "baz",
		},
		"f": []interface{}{"buz", 5},
		"g": []interface{}{},
	}, i)
	l.Remove(-1)
}

func testToFrom(t *testing.T, f func(*lua.LState, reflect.Value) lua.LValue, i interface{}, code string) {
	l := testLuaState()
	initialStackSize := l.GetTop()

	v := f(l, reflect.ValueOf(i))
	l.SetGlobal("ctx", v)
	assert.Equal(t, initialStackSize, l.GetTop())

	b := bytes.NewBufferString(code)
	fn, err := l.Load(b, "")
	require.Nil(t, err)
	l.Push(fn)
	l.Call(0, 1)
	assert.True(t, l.ToBool(-1))
	l.Remove(-1)
	assert.Equal(t, initialStackSize, l.GetTop())
}

func TestTableFromStruct(t *testing.T) {

	type Foo struct {
		A int
		B string
	}

	type Bar struct {
		C Foo
		D bool `luautil:"d"`
	}

	type Baz struct {
		Bar `luautil:",inline"`
		E   string
		F   int `luautil:"-"`
	}

	i := Baz{Bar{Foo{1, "wat"}, true}, "wut", 5}
	testToFrom(t, util.ToTableFromStruct, i, `
		if ctx.C.A ~= 1 then return false end
		if ctx.C.B ~= "wat" then return false end
		if ctx.d ~= true then return false end
		if ctx.E ~= "wut" then return false end
		if ctx.F ~= nil then return false end
		return true
	`)
}

func TestTableFromMap(t *testing.T) {
	m := map[interface{}]interface{}{
		"A": 1,
		5:   "FOO",
		true: map[string]interface{}{
			"foo": "bar",
		},
	}
	testToFrom(t, util.ToTableFromMap, m, `
		if ctx.A ~= 1 then return false end
		if ctx[5] ~= "FOO" then return false end
		if ctx[true].foo ~= "bar" then return false end
		return true
	`)
}

func TestTableFromSlice(t *testing.T) {
	s := []interface{}{
		"foo",
		true,
		4,
		[]string{
			"bar",
			"baz",
		},
	}
	testToFrom(t, util.ToTableFromSlice, s, `
		if ctx[1] ~= "foo" then return false end
		if ctx[2] ~= true then return false end
		if ctx[3] ~= 4 then return false end
		if ctx[4][1] ~= "bar" then return false end
		if ctx[4][2] ~= "baz" then return false end
		return true
	`)
}
