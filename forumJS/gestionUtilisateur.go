package forumjs

import (
	"database/sql"
	"fmt"
	"net/http"
)

func AjouterUnUtilisateur(w http.ResponseWriter, valeurEmail string, valeurMotDePasse string, nomUtilisateur string) bool {
	id := ConnecterUtilisateur(valeurEmail, valeurMotDePasse)
	if id != 0 {
		CrééUnCookie(w, id)
		return true
	}

	dsnURI := "db/user.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return false
	}

	rows, err := db.Query(`
		CREATE TABLE IF NOT EXISTS user (
			UserId INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			Email VARCHAR(80) NOT NULL,
			MotDePasse VARCHAR(80) NOT NULL,
			NomUtilisateur VARCHAR(80) NOT NULL
		);
	`)
	if err != nil {
		fmt.Println("Erreur de création :", err)
		fmt.Println(rows)
		return false
	}
	defer db.Close()

	rows, err = db.Query("SELECT UserId FROM user WHERE Email = ?;", valeurEmail)
	if err != nil {
		fmt.Println("Erreur de sélection :", err)
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			fmt.Println("scan error:", err)
			return false
		}
		fmt.Println("Cette adresse e-mail est déjà utilisé")
		return false
	}
	if err := rows.Err(); err != nil {
		fmt.Println("rows error:", err)
		return false
	}

	rows, err = db.Query(`
		INSERT INTO user (Email, MotDePasse, NomUtilisateur)
		VALUES (?, ?, ?);
	`, valeurEmail, valeurMotDePasse, nomUtilisateur)
	if err != nil {
		fmt.Println("Erreur d'insertion :", err)
		return false
	}

	iD_Utilisateur := ConnecterUtilisateur(valeurEmail, valeurMotDePasse)
	if iD_Utilisateur != 0 {
		CrééUnCookie(w, iD_Utilisateur)
		return true
	}
	return false
}

func VoirLaListeDesUtilisateurs() []StructureUtilisateur {
	liste := []StructureUtilisateur{}

	dsnURI := "db/user.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return liste
	}

	rows, err := db.Query("SELECT UserId, Email, MotDePasse, NomUtilisateur FROM user;")
	if err != nil {
		fmt.Println("Erreur de sélection :", err)
		return liste
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var email string
		var motDePasse string
		var nomUtilisateur string
		if err := rows.Scan(&id, &email, &motDePasse, &nomUtilisateur); err != nil {
			fmt.Println("scan error:", err)
			return liste
		}
		var utilisateur StructureUtilisateur
		utilisateur.iD = id
		utilisateur.email = email
		utilisateur.motDePasse = motDePasse
		utilisateur.nom = nomUtilisateur
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

func VoirUtilisateurs(id int) StructureUtilisateur {
	var utilisateur = StructureUtilisateur{}

	dsnURI := "db/user.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return utilisateur
	}

	rows, err := db.Query("SELECT UserId, Email, MotDePasse, NomUtilisateur FROM user WHERE UserId = ?;", id)
	if err != nil {
		fmt.Println("Erreur de sélection :", err)
		return utilisateur
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var email string
		var motDePasse string
		var nomUtilisateur string
		if err := rows.Scan(&id, &email, &motDePasse, &nomUtilisateur); err != nil {
			fmt.Println("scan error:", err)
			return utilisateur
		}
		utilisateur.iD = id
		utilisateur.email = email
		utilisateur.motDePasse = motDePasse
		utilisateur.nom = nomUtilisateur
		return utilisateur
	}
	if err := rows.Err(); err != nil {
		fmt.Println("rows error:", err)
	}

	return utilisateur
}

func ConnecterUtilisateur(email string, motDePasse string) int {
	dsnURI := "db/user.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return 0
	}

	rows, err := db.Query("SELECT UserId FROM user WHERE Email = ? AND MotDePasse = ?;", email, motDePasse)
	if err != nil {
		fmt.Println("Erreur de sélection :", err)
		return 0
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			fmt.Println("scan error:", err)
			return 0
		}
		return id
	}
	if err := rows.Err(); err != nil {
		fmt.Println("rows error:", err)
	}

	return 0
}
