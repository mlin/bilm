// important: go-sqlite3 must be built with: CGO_CFLAGS=-DSQLITE_ENABLE_FTS4 go build
package main

import (
	"os"
	"io"
	"bufio"
	"log"
	"syscall"
	"database/sql"
	"github.com/kless/term"
	"github.com/mattn/go-sqlite3"
)

func usage() {
	log.Fatal("Usage: cat logs | bilm_add /path/to/index")
}

func main() {
	if len(os.Args) < 2 || term.IsTerminal(syscall.Stdin) {
		usage()
	}
	dbfn := os.Args[1]
	db, err := sql.Open("sqlite3", dbfn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_ = sqlite3.ErrError

	_, err = db.Exec("create virtual table if not exists bilm using fts4(source,line_number,message)")
	if err != nil {
		log.Fatal(err)
	}

	import_source(db, "(stdin)", os.Stdin)
}

func import_source(db *sql.DB, source string, r io.Reader) {
	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt_insert, err := txn.Prepare("insert into bilm(source,line_number,message) values(?,?,?)")
	if err != nil {
		txn.Rollback()
		log.Fatal(err)
	}
	defer stmt_insert.Close()

	scanner := bufio.NewScanner(r)
	line_number := 0
	for scanner.Scan() {
		line := scanner.Text()
		line_number++
		_, err := stmt_insert.Exec(source, line_number, line)
		if err != nil {
			txn.Rollback()
			log.Fatal(err)
		}
	}

	txn.Commit()
}
