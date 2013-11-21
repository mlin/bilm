package main

import (
	"os"
	"log"
	"database/sql"
	"github.com/mattn/go-sqlite3"
)

func main() {
	os.Exit(main_impl());
}

func main_impl() int {
	if len(os.Args) != 3 {
		usage()
	}
	dbfn := os.Args[1]
	query := os.Args[2]

	db, err := sql.Open("sqlite3", dbfn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_ = sqlite3.ErrError

	var dummy string
	err = db.QueryRow("select name from sqlite_master where type='table' and name='bilm'").Scan(&dummy)
	if err == sql.ErrNoRows {
		log.Fatalf("%s does not appear to be a bilm index created using bilm_add", dbfn)
	}

	rows, err := db.Query("select message from bilm where message match ?", query)
	if err != nil {
		log.Fatal(err)
	}
	exit_code := 1
	for rows.Next() {
		var message string
		if err := rows.Scan(&message); err != nil {
			log.Fatal(err)
		}
		os.Stdout.Write([]byte(message+"\n"))
		exit_code = 0
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return exit_code
}

func usage() {
	log.Fatal("Usage: bilm_query /path/to/index \"the quick * fox\"")
}

