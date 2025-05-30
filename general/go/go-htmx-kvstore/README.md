## POCs and Pet Projects

## [go-htmx-kvstore](go-htmx-kvstore/)

Requirements:

* `gcc` *(Uses `CGO` for SQLite)*
* `go`
* [air](https://github.com/cosmtrek/air)

```
$ cd go-htmx-kvstore
$ make run
```

A Multi-tenant Key Value store built using Go, HTMX & SQLite

Topics covered:

* HTMX with Go templating ([go-htmx-kvstore/internal/web](go-htmx-kvstore/internal/web))
* Echo based HTTP Server ([go-htmx-kvstore/internal/server/server.go](https://github.com/msharran/labs/blob/main/go-htmx-kvstore/internal/server/server.go))
* SQLite DB for persisting User and KeyValues ([go-htmx-kvstore/internal/server/db.go](https://github.com/msharran/labs/blob/main/go-htmx-kvstore/internal/server/db.go))
