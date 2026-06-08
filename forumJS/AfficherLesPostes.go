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
	if len(listePostes) > 0 {
		ComplétéLaPageForum(w, r)
		AfficherPost(listePostes[0], w, r, iD_publication_commentaire == listePostes[0].Id, 0, true)

		AjouterUnCommentaire(w, r, 0, threadID, 0, true)
		if len(listePostes) > 1 {
			tableauPlacer := make([]int, 0)
			AfficherToutLesPostRécursif(w, r, &tableauPlacer, listePostes, 0, iD_publication_commentaire, 0)
		}
	} else {
		RevenirSurLaPageAccueil(w, r, 0, false, true, 0, "nouveau-sujet")
	}
}

func AfficherToutLesPostRécursif(w http.ResponseWriter, r *http.Request, tableauPlacer *[]int, listePostes []Post, answerRechercher int, iD_publication_commentaire int, décalage int) {
	for i := 1; i < len(listePostes); i++ {
		if listePostes[i].Answer == answerRechercher && !EstDansLeTableau(*tableauPlacer, listePostes[i].Id) {
			*tableauPlacer = append(*tableauPlacer, listePostes[i].Id)
			AfficherPost(listePostes[i], w, r, iD_publication_commentaire == listePostes[i].Id, décalage, false)
			AfficherToutLesPostRécursif(w, r, tableauPlacer, listePostes, listePostes[i].Id, iD_publication_commentaire, décalage+1)
		}
	}
}

/*
un poste avec answer 3 signifie qu'il est le désendant de listePostes[i+1]. Cela signifie que dans l'ordre chronologique, je doit le mettre juste après.
*/

func EstDansLeTableau(tableau []int, valeur int) bool {
	for i := 0; i < len(tableau); i++ {
		if tableau[i] == valeur {
			return true
		}
	}
	return false
}

func AfficherPost(poste Post, w http.ResponseWriter, r *http.Request, mettre_espace_commentaire bool, décalage int, premierPoste bool) {
	iD_publication := poste.Id
	iD_utilisateur_qui_poste := poste.UserId

	iD_fil_de_discussion := poste.ThreadId
	contenu_du_message := poste.Content
	date_de_publication := Date(poste.CreatedAt)
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
		"décalage":             "margin-left:" + strconv.Itoa(décalage*50) + "px;",
	}

	tmpl, err := template.ParseFiles("pages/template-post.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
		return
	}
	if premierPoste {
		tmpl, err = template.ParseFiles("pages/template-haut-file.html")
		if err != nil {
			http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
			return
		}
		listeThread, err := GetThread()
		if err != nil {
			fmt.Println(err)
		}
		id := iD_fil_de_discussion - 1
		if err == nil && id >= 0 && len(listeThread) > id {
			données["nomLabel"] = listeThread[id].Label_name
			données["nomDiscution"] = listeThread[id].Name
		}
	}

	err = tmpl.Execute(w, données)
	if err != nil {
		if isBrokenPipe(err) {
			return
		}
		fmt.Println("Erreur lors de l'exécution du template :", err)
	}

	// placer le commentaire s'il y en à un :
	if mettre_espace_commentaire {
		AjouterUnCommentaire(w, r, iD_publication, iD_fil_de_discussion, décalage, false)
	}
}

func AjouterUnCommentaire(w http.ResponseWriter, r *http.Request, iD_publication_réponce int, iD_fil_de_discussion int, décalage int, premierCommentaire bool) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	pasPosiblitéAnuler := ""
	if premierCommentaire {
		pasPosiblitéAnuler = "display:none;"
	}
	données := map[string]interface{}{
		"answer":               iD_publication_réponce,
		"iD_fil_de_discussion": iD_fil_de_discussion,
		"décalage":             "margin-left:" + strconv.Itoa(décalage*50) + "px;",
		"pasPosiblitéAnuler":   pasPosiblitéAnuler,
	}

	tmpl, err := template.ParseFiles("pages/template-commentaire.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, données)
	if err != nil {
		if isBrokenPipe(err) {
			return
		}
		fmt.Println("Erreur lors de l'exécution du template :", err)
	}
}
