package main

import (
	"os"
	"log"
	"database/sql"
	"github.com/mattn/go-sqlite3"
)

func main() {
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
	for rows.Next() {
		var message string
		if err := rows.Scan(&message); err != nil {
			log.Fatal(err)
		}
		os.Stdout.Write([]byte(message+"\n"))
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func usage() {
	log.Fatal("Usage: bilm_query /path/to/index \"the quick * fox\"")
}

