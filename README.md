bilm
====

Basic Indexed Line Matcher: index a plurality of text files to enable efficient execution of many queries for lines matching glob patterns. Uses the [SQLite3 full-text search module](http://www.sqlite.org/fts3.html).

### Usage

```{bash}
$ cat /var/log/big_log1 | bilm_add /tmp/big_logs.bilm
$ cat /var/log/big_log2 | bilm_add /tmp/big_logs.bilm
$ bilm_query /tmp/big_logs.bilm "The quick brown fox"
$ bilm_query /tmp/big_logs.bilm "jumps the * dog"
```

### Building

Assumes the [`GOPATH` environment variable](http://golang.org/doc/code.html#GOPATH) is defined. The first two lines begin by ensuring go-sqlite3 is built with `CGO_CFLAGS=-DSQLITE_ENABLE_FTS4`.

```{bash}
$ find "$GOPATH/pkg" -type f -path '*/github.com/mattn/go-sqlite3.a' -exec rm '{}' \;
$ CGO_CFLAGS=-DSQLITE_ENABLE_FTS4 go get github.com/mattn/go-sqlite3
$ go get github.com/mlin/bilm
```
