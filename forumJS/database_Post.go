package forumjs

import (
	"database/sql"
	"fmt"
	"log"
)

func InitDB() {
	dsnURI := "db/forum.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return
	}

	defer db.Close() 

	fmt.Println("Connexion à la base de donnée réussie !")

	createThreads := `
	CREATE TABLE IF NOT EXISTS Threads (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		name        TEXT NOT NULL,
		user_id     INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES utilisateurs(id)
	);`

	_, err = db.Exec(createThreads)
	if err != nil {
		fmt.Println("Erreur de création Threads :", err)
		return
	}
	fmt.Println("Table Threads créée")

	createPosts := `
	CREATE TABLE IF NOT EXISTS Posts (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id      INTEGER NOT NULL,
		thread_id    INTEGER NOT NULL,
		content      TEXT NOT NULL,
		created_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
		likes        INTEGER DEFAULT 0,
		dislikes     INTEGER DEFAULT 0,
		FOREIGN KEY (user_id) REFERENCES utilisateurs (id)
	);`

	_, err = db.Exec(createPosts)
	if err != nil {
		log.Fatal("Erreur de création Posts :", err)
		return
	}
	fmt.Println("Table Posts créée")
}
