package main 

import (
	"database/sql" 
	"fmt"
	"log" 
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB 

func InitDB() {
	var err error 
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connexion à la base de donnée réussie !")

	createThreads := `
	CREATE TABLE IF NOT EXISTS Threads (
	id          INTEGER PRIMARY KEY AUTOINCREMENT
	name        TEXT NOT NULL
	user_id     INTEGER NOT NULL
	FOREIGN KEY (user_id) REFERENCES utilusateurs(id)
	);` 

	_, err := db.Exec(createThreads) 
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Table Threads créée")

	createPosts := `
	CREATE TABLE IF NOT EXISTS Posts (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id      INTEGER NOT NULL,
	thread_id    INTEGER NOT NULL,
	content      TEXT NOT NULL,
	created_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
	likes        INTEGER DEFAULT 0
	dislikes     INTEGER DEFAULT 0,
	FOREIGN KEY (user_id) REFERENCES utilisateurs (id)
	);`


	_, err := db.Exec(createPosts) 
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("table posts créée")
} 