package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"database/sql"

	_ "modernc.org/sqlite"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"html/template"
	"strconv"
	"time"
)

type StructureUtilisateur struct {
	iD         int
	email      string
	motDePasse string
	nom        string
}

func main() {
	// Les méthode HTTP :
	http.HandleFunc("/Inscription", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// w.Write([]byte(DémarerUnePartie(informations, r)))
		email := r.FormValue("email")
		password := r.FormValue("password")
		nomUtilisateur := r.FormValue("nomUtilisateur")

		réusie := AjouterUnUtilisateur(w, email, password, nomUtilisateur)
		fmt.Println(réusie)
		if réusie {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			http.ServeFile(w, r, "inscription.html")
		}
	})

	http.HandleFunc("/Connexion", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// w.Write([]byte(DémarerUnePartie(informations, r)))
		email := r.FormValue("email")
		password := r.FormValue("password")

		réusie := false
		iD_Utilisateur := ConnecterUtilisateur(email, password)
		if iD_Utilisateur != 0 {
			CrééUnCookie(w, iD_Utilisateur)
			réusie = true
		}
		fmt.Println(réusie)

		if réusie {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			http.ServeFile(w, r, "inscription.html")
		}
	})

	http.HandleFunc("/PageInscription", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, "inscription.html")
	})

	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("./style"))))

	// Au démarage du serveur :
	log.Println("Serveur lancé sur http://localhost:8080")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		idUtilisateur := VérifierCookie(r)
		nomAAfficher := "Invité"

		if idUtilisateur != 0 {
			utilisateur := VoirUtilisateurs(idUtilisateur)
			if utilisateur.nom != "" {
				nomAAfficher = utilisateur.nom
			}
		}

		données := map[string]interface{}{
			"NomUtilisateur": nomAAfficher,
		}

		tmpl, err := template.ParseFiles("main.html")
		if err != nil {
			http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, données)
		if err != nil {
			fmt.Println("Erreur lors de l'exécution du template :", err)
		}
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

var cleSecrete = []byte("shjzaqfjkffzf5ver6ezrcf8rceez569z")

func CrééUnCookie(w http.ResponseWriter, id int) {
	durée := 3 * time.Hour
	expiration := time.Now().Add(durée)

	idTexte := strconv.Itoa(id)

	mac := hmac.New(sha256.New, cleSecrete)
	mac.Write([]byte(idTexte))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	valeurSecurisée := idTexte + "." + signature

	cookie := &http.Cookie{
		Name:     "session_utilisateur",
		Value:    valeurSecurisée,
		Expires:  expiration,
		MaxAge:   int(durée.Seconds()),
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	}

	http.SetCookie(w, cookie)
}

func VérifierCookie(r *http.Request) int {
	cookie, err := r.Cookie("session_utilisateur")
	if err != nil {
		return 0
	}

	valeur := cookie.Value

	var idTexte, signatureRecue string
	for i := 0; i < len(valeur); i++ {
		if valeur[i] == '.' {
			idTexte = valeur[:i]
			signatureRecue = valeur[i+1:]
			break
		}
	}

	if idTexte == "" || signatureRecue == "" {
		return 0
	}

	mac := hmac.New(sha256.New, cleSecrete)
	mac.Write([]byte(idTexte))
	signatureAttendue := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if signatureRecue != signatureAttendue {
		fmt.Println("Tentative de modification de cookie détectée !")
		return 0
	}

	id, err := strconv.Atoi(idTexte)
	if err != nil {
		return 0
	}

	return id
}
