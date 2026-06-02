package forumjs

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func AfficherToutLesPost(threadID int, w http.ResponseWriter, r *http.Request, iD_publication_commentaire int) {
	// -1 si il n'y a rien.
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
		AfficherPost(listePostes[i], w, r, iD_publication_commentaire-1 == i)
	}
}

func AfficherPost(poste Post, w http.ResponseWriter, r *http.Request, mettre_espace_commentaire bool) {
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

	iconeAime := "images/aime.svg"
	iconeAimePas := "images/aime.svg"

	idUtilisateur := VérifierCookie(r)
	if idUtilisateur != 0 {
		if LireTableauInteractionUtilisateur(w, r, idUtilisateur, iD_publication, iD_fil_de_discussion, "likes") {
			iconeAime = "images/aimeActif.svg"
		}
		if LireTableauInteractionUtilisateur(w, r, idUtilisateur, iD_publication, iD_fil_de_discussion, "dislikes") {
			iconeAimePas = "images/aimeActif.svg"
		}
	}

	données := map[string]interface{}{
		"nom_utilisateur":      nom_utilisateur,
		"contenu_du_message":   contenu_du_message,
		"date_de_publication":  date_de_publication,
		"nombre_de_aime":       nombre_de_aime,
		"nombre_de_aime_pas":   nombre_de_aime_pas,
		"iD_publication":       iD_publication,
		"iD_fil_de_discussion": iD_fil_de_discussion,
		"iconeAime":            iconeAime,
		"iconeAimePas":         iconeAimePas,
		"nomPosteID":           "post-" + strconv.Itoa(iD_publication),
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

	// placer le commentaire s'il y en à un :
	if mettre_espace_commentaire {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		données := map[string]interface{}{
			"nom_utilisateur":      nom_utilisateur,
			"contenu_du_message":   "Réponce !",
			"date_de_publication":  date_de_publication,
			"nombre_de_aime":       nombre_de_aime,
			"nombre_de_aime_pas":   nombre_de_aime_pas,
			"iD_publication":       iD_publication,
			"iD_fil_de_discussion": iD_fil_de_discussion,
			"iconeAime":            iconeAime,
			"iconeAimePas":         iconeAimePas,
			"nomPosteID":           "post-" + strconv.Itoa(iD_publication),
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
}
