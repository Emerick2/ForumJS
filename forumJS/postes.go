package forumjs

import (
	"fmt"
	"net/http"
	"text/template"
)

func AfficherPost(id int, w http.ResponseWriter, r *http.Request) {
	iD_publication := 1
	iD_utilisateur_qui_poste := 1

	iD_fil_de_discussion := 3
	contenu_du_message := "blabla"
	date_de_publication := "29 mai 2026"
	nombre_de_aime := 15
	nombre_de_aime_pas := 2

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
