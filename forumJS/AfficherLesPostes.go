package forumjs

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func AfficherToutLesPost(threadID int, w http.ResponseWriter, r *http.Request, iD_publication_commentaire int) {
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

	premierIDPoste := listePostes[0].Id

	listePostes = AjouterDonnéesPostes(listePostes, w, r, iD_publication_commentaire, true, premierIDPoste)

	var tableauPlacer []int
	var nouvelleListe []Post
	nouvelleListe = append(nouvelleListe, listePostes[0])
	AfficherToutLesPostRécursif(w, r, &tableauPlacer, listePostes, premierIDPoste, &nouvelleListe)

	listePostes = nouvelleListe

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("pages/discution.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
		return
	}

	données := map[string]interface{}{
		"ListePostes": listePostes,
	}

	err = tmpl.Execute(w, données)
	if err != nil {
		if isBrokenPipe(err) {
			return
		}
		fmt.Println("Erreur lors de l'exécution du template :", err)
	}
}

func AjouterDonnéesPostes(listePostes []Post, w http.ResponseWriter, r *http.Request, iD_publication_commentaire int, utiliserInterfacePublication bool, premierIDPoste int) []Post {
	for i := 0; i < len(listePostes); i++ {
		unPost := listePostes[i]
		unPost.NameUser = "Compte suprimé"
		valeur := VoirUtilisateurs(unPost.UserId)
		if valeur.nom != "" {
			unPost.NameUser = valeur.nom
		}

		unPost.CreatedAtText = Date(unPost.CreatedAt)

		unPost.IconeLike = "/images/aime.svg"
		unPost.IconeDislike = "/images/aime.svg"

		idUtilisateur := VérifierCookie(r)
		if idUtilisateur != 0 {
			if LireTableauInteractionUtilisateur(w, r, idUtilisateur, unPost.Id, unPost.ThreadId, "likes") {
				unPost.IconeLike = "/images/aimeActif.svg"
			}
			if LireTableauInteractionUtilisateur(w, r, idUtilisateur, unPost.Id, unPost.ThreadId, "dislikes") {
				unPost.IconeDislike = "/images/aimeActif.svg"
			}
		}

		unPost.NameOfTheIdPost = "post-" + strconv.Itoa(unPost.Id)

		if (unPost.Answer != 0 && unPost.Answer != premierIDPoste) && utiliserInterfacePublication {
			unPost.TheMargin = "margin-left:50px;"
			unPost.BlockComments = "display:none;"
		}

		if utiliserInterfacePublication {
			if i != 0 {
				unPost.BlockShare = "display:none;"
			}
			if i != 0 && iD_publication_commentaire != unPost.Id {
				unPost.BlockNewComments = "display:none;"
			} else {
				unPost.BlockComments = "display:none;"
			}
			if i == 0 {
				unPost.OptionToCancel = "display:none;"
				// données du fil de discution :
				dsnURI := "db/threads.db"
				db, err := sql.Open("sqlite", dsnURI)
				if err != nil {
					fmt.Println("Erreur d'ouverture :", err)
				}
				defer db.Close()

				requete := fmt.Sprintf("SELECT name, label_name FROM Threads WHERE id = ?")

				err = db.QueryRow(requete, unPost.ThreadId).Scan(
					&unPost.NameThread,
					&unPost.LabelThread,
				)

				unPost.TexteFil = unPost.NameThread + " [" + unPost.LabelThread + "]"
			}
		} else {
			unPost.BlockNewComments = "display:none;"
		}
		listePostes[i] = unPost
	}
	return listePostes
}

func AfficherToutLesPostRécursif(w http.ResponseWriter, r *http.Request, tableauPlacer *[]int, listePostes []Post, answerRechercher int, nouvelleListe *[]Post) {
	for i := 1; i < len(listePostes); i++ {
		if listePostes[i].Answer == answerRechercher && !EstDansLeTableau(*tableauPlacer, listePostes[i].Id) {
			*tableauPlacer = append(*tableauPlacer, listePostes[i].Id)
			*nouvelleListe = append(*nouvelleListe, listePostes[i])
			AfficherToutLesPostRécursif(w, r, tableauPlacer, listePostes, listePostes[i].Id, nouvelleListe)
		}
	}
}

func EstDansLeTableau(tableau []int, valeur int) bool {
	for i := 0; i < len(tableau); i++ {
		if tableau[i] == valeur {
			return true
		}
	}
	return false
}
