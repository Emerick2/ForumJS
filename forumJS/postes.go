package forumjs

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func AfficherToutLesPost(threadID int, w http.ResponseWriter, r *http.Request) {
	dsnURI := "db/forum.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return
	}

	defer db.Close()

	listePostes, err := GetPostsByThread(threadID, db)
	if err != nil {
		fmt.Println("Erreur lors de la récupération des posts :", err)
		return
	}

	for i := 0; i < len(listePostes); i++ {
		AfficherPost(listePostes[i], w, r)
	}
}

func AfficherPost(poste Post, w http.ResponseWriter, r *http.Request) {
	iD_publication := poste.Id
	iD_utilisateur_qui_poste := poste.UserId

	iD_fil_de_discussion := poste.ThreadId
	contenu_du_message := poste.Content
	date_de_publication := poste.CreatedAt
	nombre_de_aime := poste.Likes
	nombre_de_aime_pas := poste.Dislikes

	nom_utilisateur := "Compte suprimé"
	valeur := VoirUtilisateurs(iD_utilisateur_qui_poste)
	if valeur.nom != "" {
		nom_utilisateur = valeur.nom
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	données := map[string]interface{}{
		"nom_utilisateur":      nom_utilisateur,
		"contenu_du_message":   contenu_du_message,
		"date_de_publication":  date_de_publication,
		"nombre_de_aime":       nombre_de_aime,
		"nombre_de_aime_pas":   nombre_de_aime_pas,
		"iD_publication":       iD_publication,
		"iD_fil_de_discussion": iD_fil_de_discussion,
	}

	tmpl, err := template.ParseFiles("pages/template-post.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, données)
	if err != nil {
		fmt.Println("Erreur lors de l'exécution du template :", err)
	}
}

func InteractionPost(w http.ResponseWriter, r *http.Request) {
	nomAction := r.FormValue("nomAction")
	valeur := (r.FormValue("iD_publication"))
	iD_publication, err := strconv.Atoi(valeur)
	if err != nil {
		fmt.Println(err)
	}
	valeur = r.FormValue("iD_fil_de_discussion")
	iD_fil_de_discussion, err := strconv.Atoi(valeur)
	if err != nil {
		fmt.Println(err)
	}
	idUtilisateur := VérifierCookie(r)
	if idUtilisateur == 0 {
		return
	}

	dsnURI := "db/forum.db"
	// dsnURIUtilisateur := "db/user.db"
	if nomAction == "aime" {
		changement := 1
		if LireTableauInteractionUtilisateur(w, r, idUtilisateur, iD_publication, iD_fil_de_discussion, "likes") {
			changement = -1
		}
		if LireTableauInteractionUtilisateur(w, r, idUtilisateur, iD_publication, iD_fil_de_discussion, "dislikes") {
			SauvegarderUneValeur(w, r, dsnURI, iD_publication, iD_fil_de_discussion, "dislikes", -1, "Posts")
			SauvegarderTableauInteractionUtilisateur(w, r, idUtilisateur, iD_publication, iD_fil_de_discussion, "dislikes", -1)
		}
		SauvegarderUneValeur(w, r, dsnURI, iD_publication, iD_fil_de_discussion, "likes", changement, "Posts")
		SauvegarderTableauInteractionUtilisateur(w, r, idUtilisateur, iD_publication, iD_fil_de_discussion, "likes", changement)
	} else if nomAction == "aimePas" {
		changement := 1
		if LireTableauInteractionUtilisateur(w, r, idUtilisateur, iD_publication, iD_fil_de_discussion, "dislikes") {
			changement = -1
		}
		if LireTableauInteractionUtilisateur(w, r, idUtilisateur, iD_publication, iD_fil_de_discussion, "likes") {
			SauvegarderUneValeur(w, r, dsnURI, iD_publication, iD_fil_de_discussion, "likes", -1, "Posts")
			SauvegarderTableauInteractionUtilisateur(w, r, idUtilisateur, iD_publication, iD_fil_de_discussion, "likes", -1)
		}
		SauvegarderUneValeur(w, r, dsnURI, iD_publication, iD_fil_de_discussion, "dislikes", changement, "Posts")
		SauvegarderTableauInteractionUtilisateur(w, r, idUtilisateur, iD_publication, iD_fil_de_discussion, "dislikes", changement)
	}
}

func SauvegarderUneValeur(w http.ResponseWriter, r *http.Request, dsnURI string, iD_publication int, iD_fil_de_discussion int, clef string, modification int, nomTable string) {
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return
	}

	defer db.Close()

	if clef != "likes" && clef != "dislikes" {
		http.Error(w, "Colonne invalide", http.StatusBadRequest)
		return
	}

	requete := fmt.Sprintf("SELECT %s FROM %s WHERE id = ? AND thread_id = ? LIMIT 1", clef, nomTable)
	var valeurRecup int
	err = db.QueryRow(requete, iD_publication, iD_fil_de_discussion).Scan(&valeurRecup)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Poste non trouvé", http.StatusNotFound)
			return
		}
		http.Error(w, "Erreur lors de la récupération des données", http.StatusInternalServerError)
		fmt.Println("QueryRow error:", err)
		return
	}

	valeurObtenu := valeurRecup + modification
	updateReq := fmt.Sprintf("UPDATE %s SET %s = ? WHERE id = ? AND thread_id = ?", nomTable, clef)
	_, err = db.Exec(updateReq, valeurObtenu, iD_publication, iD_fil_de_discussion)
	if err != nil {
		http.Error(w, "Erreur lors de la sauvegarde des données", http.StatusInternalServerError)
		fmt.Println("Exec update error:", err)
		return
	}
}

func LireUneValeur(w http.ResponseWriter, r *http.Request, dsnURI string, iD_publication int, iD_fil_de_discussion int, clef string, modification int, nomTable string) int {
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return 0
	}

	defer db.Close()

	if clef != "likes" && clef != "dislikes" {
		http.Error(w, "Colonne invalide", http.StatusBadRequest)
		return 0
	}

	requete := fmt.Sprintf("SELECT %s FROM %s WHERE id = ? AND thread_id = ? LIMIT 1", clef, nomTable)
	var valeurRecup int
	err = db.QueryRow(requete, iD_publication, iD_fil_de_discussion).Scan(&valeurRecup)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Poste non trouvé", http.StatusNotFound)
			return 0
		}
		http.Error(w, "Erreur lors de la récupération des données", http.StatusInternalServerError)
		fmt.Println("QueryRow error:", err)
		return 0
	}

	return valeurRecup + modification
}

func LireTableauInteractionUtilisateur(w http.ResponseWriter, r *http.Request, UserId int, iD_publication int, iD_fil_de_discussion int, clef string) bool {
	// retourne true si l'utilisateur a déjà interagi pour ce post.

	dsnURI := "db/interactionUtilisateur.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return false
	}

	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS interactionUtilisateur (
			Id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			UserId INTEGER NOT NULL,
			likes INTEGER NOT NULL DEFAULT 0,
			dislikes INTEGER NOT NULL DEFAULT 0,
			iD_publication INTEGER NOT NULL,
			iD_fil_de_discussion INTEGER NOT NULL
		);
	`)
	if err != nil {
		fmt.Println("Erreur de création Threads :", err)
		return false
	}

	requete := fmt.Sprintf("SELECT %s FROM interactionUtilisateur WHERE UserId = ? AND iD_publication = ? AND iD_fil_de_discussion = ? LIMIT 1", clef)
	var valeurRecup int
	err = db.QueryRow(requete, UserId, iD_publication, iD_fil_de_discussion).Scan(&valeurRecup)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("ErrNoRows")
			return false
		}
		fmt.Println("QueryRow error:", err)
		return false
	}

	return valeurRecup > 0
}

func SauvegarderTableauInteractionUtilisateur(w http.ResponseWriter, r *http.Request, UserId int, iD_publication int, iD_fil_de_discussion int, clef string, nouvelValeur int) {
	dsnURI := "db/interactionUtilisateur.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return
	}

	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS interactionUtilisateur (
			Id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			UserId INTEGER NOT NULL,
			likes INTEGER NOT NULL DEFAULT 0,
			dislikes INTEGER NOT NULL DEFAULT 0,
			iD_publication INTEGER NOT NULL,
			iD_fil_de_discussion INTEGER NOT NULL
		);
	`)
	if err != nil {
		fmt.Println("Erreur de création Threads :", err)
		return
	}

	_, err = db.Exec(`
		INSERT OR IGNORE INTO interactionUtilisateur (UserId, likes, dislikes, iD_publication, iD_fil_de_discussion)
		VALUES (?, 0, 0, ?, ?)
	`, UserId, iD_publication, iD_fil_de_discussion)
	if err != nil {
		fmt.Println("Erreur d'insertion interactionUtilisateur :", err)
		return
	}

	updateReq := fmt.Sprintf("UPDATE interactionUtilisateur SET %s = ? WHERE UserId = ? AND iD_publication = ? AND iD_fil_de_discussion = ?", clef)
	_, err = db.Exec(updateReq, nouvelValeur, UserId, iD_publication, iD_fil_de_discussion)
	if err != nil {
		http.Error(w, "Erreur lors de la sauvegarde des données ici", http.StatusInternalServerError)
		fmt.Println("Exec update error:", err)
		return
	}
}
