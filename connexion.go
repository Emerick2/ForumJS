package main

import (
	"fmt"

	"database/sql"

	_ "modernc.org/sqlite"
)

type StructureUtilisateur struct {
	iD         int
	email      string
	motDePasse string
}

func main() {
	if true {
		email := ""
		fmt.Print("Quel est votre e-mail ? ")
		fmt.Scan(&email)

		motDePasse := ""
		fmt.Print("Quel est votre mot de passe ? ")
		fmt.Scan(&motDePasse)

		AjouterUnUtilisateur(email, motDePasse)
	}
	if false {
		fmt.Println(VoirUtilisateurs(5))
	}
}

func AjouterUnUtilisateur(valeurEmail string, valeurMotDePasse string) {
	dsnURI := "db/user.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return
	}

	rows, err := db.Query(`
		CREATE TABLE IF NOT EXISTS user (
			UserId INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			Email VARCHAR(80) NOT NULL,
			MotDePasse VARCHAR(80) NOT NULL
		);
	`)
	if err != nil {
		fmt.Println("Erreur de création :", err)
		fmt.Println(rows)
		return
	}
	defer db.Close()

	rows, err = db.Query(`
		INSERT INTO user (Email, MotDePasse)
		VALUES (?, ?);
	`, valeurEmail, valeurMotDePasse)
	if err != nil {
		fmt.Println("Erreur d'insertion :", err)
		return
	}

	VoirLaListeDesUtilisateurs()
}

func VoirLaListeDesUtilisateurs() []StructureUtilisateur {
	liste := []StructureUtilisateur{}

	dsnURI := "db/user.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return liste
	}

	rows, err := db.Query("SELECT UserId, Email, MotDePasse FROM user;")
	if err != nil {
		fmt.Println("Erreur de sélection :", err)
		return liste
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var email string
		var motDePasse string
		if err := rows.Scan(&id, &email, &motDePasse); err != nil {
			fmt.Println("scan error:", err)
			return liste
		}
		var utilisateur StructureUtilisateur
		utilisateur.iD = id
		utilisateur.email = email
		utilisateur.motDePasse = motDePasse
		liste = append(liste, utilisateur)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("rows error:", err)
	}

	for i := 0; i < len(liste); i++ {
		fmt.Println(liste[i])
	}
	return liste
}

func VoirUtilisateurs(id int) StructureUtilisateur{
	var utilisateur = StructureUtilisateur{}

	dsnURI := "db/user.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return utilisateur
	}

	rows, err := db.Query("SELECT UserId, Email, MotDePasse FROM user WHERE UserId = ?;",id)
	if err != nil {
		fmt.Println("Erreur de sélection :", err)
		return utilisateur
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var email string
		var motDePasse string
		if err := rows.Scan(&id, &email, &motDePasse); err != nil {
			fmt.Println("scan error:", err)
			return utilisateur
		}
		utilisateur.iD = id
		utilisateur.email = email
		utilisateur.motDePasse = motDePasse
		return utilisateur
	}
	if err := rows.Err(); err != nil {
		fmt.Println("rows error:", err)
	}

	return utilisateur
}