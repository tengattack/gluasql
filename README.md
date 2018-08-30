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

// Preload LuaSocket modules
gluasocket.Preload(L)
```

### Query

```lua
mysql = require('mysql')
c = mysql.new()
ok, err = c:connect({ host = '127.0.0.1', port = 3306, database = 'test', user = 'user', password = 'pass' })
if ok then
  res, err = c:query('SELECT * FROM mytable LIMIT 2')
  dump(res)
end
```

## Testing

```bash
$ go test -coverprofile=/tmp/go-code-cover github.com/tengattack/gluasql...
?       github.com/tengattack/gluasql   [no test files]
ok      github.com/tengattack/gluasql/mysql     1.116s  coverage: 69.5% of statements
ok      github.com/tengattack/gluasql/util      0.070s  coverage: 76.9% of statements
```

## License

MIT
