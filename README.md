# LuaSql for GopherLua

A native Go implementation of SQL client for the [GopherLua](https://github.com/yuin/gopher-lua) VM.

## Using

### Loading Modules

```go
import (
	"github.com/tengattack/gluasql"
)

// Bring up a GopherLua VM
L := lua.NewState()
defer L.Close()

// Preload LuaSql modules
gluasql.Preload(L)
```

Or if we only need `mysql` module:

```go
import (
	mysql "github.com/tengattack/gluasql/mysql"
)

// Bring up a GopherLua VM
L := lua.NewState()
defer L.Close()

L.PreloadModule("mysql", mysql.Loader)
```

### MySQL

```lua
mysql = require('mysql')
c = mysql.new()
ok, err = c:connect({ host = '127.0.0.1', port = 3306, database = 'test', user = 'user', password = 'pass' })
if ok then
  res, err = c:query('SELECT * FROM mytable LIMIT 2')
  dump(res)
end
```

### PostgreSQL

import [lib/pq](https://github.com/lib/pq) first:

```go
import (
	_ "github.com/lib/pq"
)
```

```lua
mysql = require('mysql')
c = mysql.new()
ok, err = c:connect('postgres://host=127.0.0.1 port=5432 user=user password=123456 dbname=foo sslmode=disable')
if ok then
  res, err = c:query('SELECT * FROM mytable LIMIT 2')
  dump(res)
end
```

### SQLite

Since it depends `go-sqlite3`, we need `gcc` to compile SQLite module, more details:
[https://github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)

```lua
sqlite3 = require('sqlite3')
c = sqlite3.new()
ok, err = c:open('test.db', { cache = 'shared' })
if ok then
  res, err = c:query('SELECT * FROM mytable LIMIT 2')
  dump(res)
end
```

## Testing

```bash
$ go test -coverprofile=/tmp/go-code-cover github.com/tengattack/gluasql...
ok      github.com/tengattack/gluasql/mysql     0.020s  coverage: 79.9% of statements
ok      github.com/tengattack/gluasql/sqlite3   0.019s  coverage: 71.3% of statements
ok      github.com/tengattack/gluasql/util      0.015s  coverage: 76.9% of statements
```

## License

MIT
