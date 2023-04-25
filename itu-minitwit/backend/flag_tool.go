package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const docStr = `ITU-Minitwit Tweet Flagging Tool

Usage:
  flag_tool <tweet_id>...
  flag_tool -i
  flag_tool -h
Options:
-h            Show this screen.
-i            Dump all tweets and authors to STDOUT.
`

func callback(rows *sql.Rows) {
	var col0, col1, col2, col4 string
	for rows.Next() {
		err := rows.Scan(&col0, &col1, &col2, &col4)
		if err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Println(col0, col1, col2, col4)
	}
}

func errorCheck(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {
	args := os.Args[1:]

	db, err := sql.Open("sqlite3", "./../tmp/minitwit.db")
	errorCheck(err)
	defer db.Close()

	if len(args) == 1 && args[0] == "-h" {
		fmt.Print(docStr)
		return
	}
	if len(args) == 1 && args[0] == "-i" {
		rows, err := db.Query("SELECT * FROM message")
		errorCheck(err)
		callback(rows)
	}
	if args[0] != "-i" && args[0] != "-h" {
		fmt.Println("Flagging entries: ")
		for i := 0; i < len(args); i++ {
			_, err = db.Exec("UPDATE message SET flagged=1 WHERE message_id=" + args[i])
			if err != nil {
				fmt.Println("Error executing query:", err)
				return
			}
			fmt.Println("Flagged entry:", args[i])
		}
	}
}
