package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// connect to db 
// create db tables if they don't exist

//path to db
dbPath := "./../../tmp/minitwit.db"

func connect_db(){
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.println("Connected to db")


