package gluasql_mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	util "github.com/tengattack/gluasql/util"
	lua "github.com/yuin/gopher-lua"
)

// errors
var ErrConnectionString = errors.New("options or connection string excepted")

// ParseConnectionString parse options or connection string to golang sql driverName and dsn
func ParseConnectionString(cs interface{}, timeout time.Duration) (driverName string, dsn string, err error) {
	if cs == nil {
		return "", "", ErrConnectionString
	} else if connectionString, ok := cs.(string); ok {
		if connectionString == "" {
			return "", "", ErrConnectionString
		}
		driverName = "mysql"
		pos := strings.Index(connectionString, "://")
		if pos >= 0 {
			driverName = connectionString[:pos]
			dsn = connectionString[pos+3:]
		} else {
			dsn = connectionString
		}
		if timeout > 0 {
			switch driverName {
			case "mysql":
				// https://github.com/go-sql-driver/mysql
				stimeout := timeout.String()
				if !strings.Contains(dsn, "?") {
					dsn += "?"
				} else {
					dsn += "&"
				}
				dsn += fmt.Sprintf("readTimeout=%s&writeTimeout=%s", stimeout, stimeout)
			case "postgres":
				// https://github.com/lib/pq
				dsn += " connect_timeout=" + strconv.Itoa(int(timeout.Seconds()))
			}
			// TODO: sqlserver etc.
		}
		return
	} else if options, ok := cs.(map[string]interface{}); ok {
		driverName = "mysql"
		host, _ := options["host"].(string)
		if host == "" {
			host = "127.0.0.1"
		}
		port, _ := options["port"].(int)
		if port == 0 {
			port = 3306
		}
		database, _ := options["database"].(string)
		user, _ := options["user"].(string)
		password, _ := options["password"].(string)
		charset, _ := options["charset"].(string)

		// current support tcp connection only
		dsn = fmt.Sprintf("tcp(%s:%d)/%s", host, port, database)
		if user != "" {
			if password != "" {
				dsn = fmt.Sprintf("%s:%s@", user, password) + dsn
			} else {
				dsn = fmt.Sprintf("%s@", user) + dsn
			}
		}

		query := url.Values{}
		if charset != "" {
			query.Set("charset", charset)
		}
		if timeout > 0 {
			stimeout := timeout.String()
			query.Set("readTimeout", stimeout)
			query.Set("writeTimeout", stimeout)
		}

		s := query.Encode()
		if s != "" {
			dsn += "?" + s
		}
		return
	}
	return "", "", ErrConnectionString
}

func clientConnectMethod(L *lua.LState) int {
	client := checkClient(L)
	cs := util.GetValue(L, 2)

	driverName, dsn, err := ParseConnectionString(cs, client.Timeout)
	if err != nil {
		L.ArgError(2, err.Error())
		return 0
	}

	client.DB, err = sql.Open(driverName, dsn)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}
