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
	iD_publication := r.FormValue("iD_publication")
	iD_fil_de_discussion := r.FormValue("iD_fil_de_discussion")

	fmt.Println(nomAction)
	fmt.Println(iD_publication)
	fmt.Println(iD_fil_de_discussion)

	dsnURI := "db/forum.db"
	if nomAction == "aime" {
		SauvegarderUneValeur(w, r, dsnURI, iD_publication, iD_fil_de_discussion, "likes", 1)
	} else if nomAction == "aimePas" {
		SauvegarderUneValeur(w, r, dsnURI, iD_publication, iD_fil_de_discussion, "dislikes", 1)
	}
}

func SauvegarderUneValeur(w http.ResponseWriter, r *http.Request, dsnURI string, iD_publication string, iD_fil_de_discussion string, clef string, modification int) {
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

	idPub, err := strconv.Atoi(iD_publication)
	if err != nil {
		http.Error(w, "ID de publication invalide", http.StatusBadRequest)
		return
	}
	idThread, err := strconv.Atoi(iD_fil_de_discussion)
	if err != nil {
		http.Error(w, "ID de fil invalide", http.StatusBadRequest)
		return
	}

	requete := fmt.Sprintf("SELECT %s FROM Posts WHERE id = ? AND thread_id = ? LIMIT 1", clef)
	var valeurRecup int
	err = db.QueryRow(requete, idPub, idThread).Scan(&valeurRecup)
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

	updateReq := fmt.Sprintf("UPDATE Posts SET %s = ? WHERE id = ? AND thread_id = ?", clef)
	_, err = db.Exec(updateReq, valeurObtenu, idPub, idThread)
	if err != nil {
		http.Error(w, "Erreur lors de la sauvegarde des données", http.StatusInternalServerError)
		fmt.Println("Exec update error:", err)
		return
	}
}
