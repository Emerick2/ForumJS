package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"database/sql"

	_ "modernc.org/sqlite"
)

type StructureUtilisateur struct {
	iD         int
	email      string
	motDePasse string
	nom        string
}

func main() {
	// if true {
	// 	email := ""
	// 	fmt.Print("Quel est votre e-mail ? ")
	// 	fmt.Scan(&email)

	// 	motDePasse := ""
	// 	fmt.Print("Quel est votre mot de passe ? ")
	// 	fmt.Scan(&motDePasse)

	// 	AjouterUnUtilisateur(email, motDePasse)
	// }
	// if false {
	// 	fmt.Println(VoirUtilisateurs(5))
	// }

	// Les méthode HTTP :
	http.HandleFunc("/Inscription", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// w.Write([]byte(DémarerUnePartie(informations, r)))
		email := r.FormValue("email")
		password := r.FormValue("password")

		réusie := AjouterUnUtilisateur(email, password)
		fmt.Println(réusie)
	})

	http.HandleFunc("/Connexion", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// w.Write([]byte(DémarerUnePartie(informations, r)))
		email := r.FormValue("email")
		password := r.FormValue("password")

		réusie := false
		iD_Utilisateur := ConnecterUtilisateur(email, password)
		if iD_Utilisateur != 0 {
			CrééUnCookie(iD_Utilisateur)
			réusie = true
		}
		fmt.Println(réusie)
	})

	// Au démarage du serveur :
	log.Println("Serveur lancé sur http://localhost:8080")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, "main.html")
	})

	http.HandleFunc("/open", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		go func() {
			_ = exec.Command("xdg-open", "http://localhost:8080/").Start()
		}()
		w.Write([]byte("Attempted to open browser"))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}

// Les autres fonctions :

func AjouterUnUtilisateur(valeurEmail string, valeurMotDePasse string) bool {
	id := ConnecterUtilisateur(valeurEmail, valeurMotDePasse)
	if id != 0{
		CrééUnCookie(id)
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
			MotDePasse VARCHAR(80) NOT NULL
		);
	`)
	if err != nil {
		fmt.Println("Erreur de création :", err)
		fmt.Println(rows)
		return false
	}
	defer db.Close()

	rows, err = db.Query(`
		INSERT INTO user (Email, MotDePasse)
		VALUES (?, ?);
	`, valeurEmail, valeurMotDePasse)
	if err != nil {
		fmt.Println("Erreur d'insertion :", err)
		return false
	}

	iD_Utilisateur := ConnecterUtilisateur(valeurEmail, valeurMotDePasse)
	if iD_Utilisateur != 0 {
		CrééUnCookie(iD_Utilisateur)
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

func VoirUtilisateurs(id int) StructureUtilisateur {
	var utilisateur = StructureUtilisateur{}

	dsnURI := "db/user.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return utilisateur
	}

	rows, err := db.Query("SELECT UserId, Email, MotDePasse FROM user WHERE UserId = ?;", id)
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

func CrééUnCookie(id int) {
	// Les cookies dures 3h
	// duréeDeVieDuCookie
}
