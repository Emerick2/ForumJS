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

	motDePasse := ""
	fmt.Print("Quel est votre mot de passe ? ")
	fmt.Scan(&motDePasse)


	AjouterUnUtilisateur(email, motDePasse)
}

func AjouterUnUtilisateur(valeurEmail string, valeurMotDePasse string) {
	dsnURI := "/tmp/testbase.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil { fmt.Println("Erreur d'ouverture :", err); return }

	rows, err := db.Query(`
		CREATE TABLE IF NOT EXISTS user (
			UserId INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			Email VARCHAR(80) NOT NULL,
			MotDePasse VARCHAR(80) NOT NULL
		);
	`)
	if err != nil { fmt.Println("Erreur de création :", err); return }
	defer db.Close()

	rows, err = db.Query(`
		INSERT INTO user (Email, MotDePasse)
		VALUES (?, ?);
	`, valeurEmail, valeurMotDePasse)
	if err != nil { fmt.Println("Erreur d'insertion :", err); return }


	rows, err = db.Query("SELECT UserId, Email, MotDePasse FROM user;")
    if err != nil { fmt.Println("Erreur de sélection :", err); return }
    defer rows.Close()

    for rows.Next() {
      var id int
      var email string
	  var motDePasse string
      if err := rows.Scan(&id, &email, &motDePasse); err != nil {
        fmt.Println("scan error:", err)
        return
      }
      fmt.Println(id, email, motDePasse)
    }
    if err := rows.Err(); err != nil {
      fmt.Println("rows error:", err)
    }
}
