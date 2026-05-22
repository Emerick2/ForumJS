package main

import (
	"fmt"

	"database/sql"

	_ "modernc.org/sqlite"
)

func main() {
	email := ""
	fmt.Print("Quel est votre e-mail ? ")
	fmt.Scan(&email)

	AjouterDansLaBaseDeDonnée("utilisateur", "email", email)
}

func AjouterDansLaBaseDeDonnée(nomDeLaTable string, nomDeLaClef string, valeur string) {
	dsnURI := "/tmp/testbase.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil { fmt.Println("open error:", err); return }

	rows, err := db.Query(`
		CREATE TABLE IF NOT EXISTS user (
			UserId INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			Email VARCHAR(80),
			MotDePasse VARCHAR(80)
		);
	`)
	// rows, err := db.Query(`
	// 	CREATE TABLE IF NOT EXISTS user (
	// 		UserId INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	// 		Email VARCHAR(80) NOT NULL,
	// 		MotDePasse VARCHAR(80) NOT NULL
	// 	);
	// `)
	if err != nil { fmt.Println("create error:", err); return }
	defer db.Close()

	rows, err = db.Query(`
		INSERT INTO user (Email)
		VALUES (?);
	`, valeur)
	if err != nil { fmt.Println("insert error:", err); return }


	rows, err = db.Query("SELECT UserId, Email FROM user;")
    if err != nil { fmt.Println("select error:", err); return }
    defer rows.Close()

    for rows.Next() {
      var id int
      var email string
      if err := rows.Scan(&id, &email); err != nil {
        fmt.Println("scan error:", err)
        return
      }
      fmt.Println(id, email)
    }
    if err := rows.Err(); err != nil {
      fmt.Println("rows error:", err)
    }
}
