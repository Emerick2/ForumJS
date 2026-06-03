package handlers

import (
	"net/http"

	"forumJS/internal/models"
)

// handleVoirPost affiche un post avec ses commentaires
// ticket #3 — Page du post
func (app *App) handleVoirPost(w http.ResponseWriter, r *http.Request) {
	// on récupère l'id du post depuis l'URL (/post/{id})
	idPost := r.PathValue("id")

	utilisateur := app.chercherUtilisateurSession(r)

	idUtilisateur := idOuVide(utilisateur)

	// on cherche le post dans la base de données
	post, err := app.db.GetPostByID(idPost, idUtilisateur)
	if err != nil {
		app.afficherErreur(w, r, 500, "Erreur serveur")
		return
	}
	if post == nil {
		app.afficherErreur(w, r, 404, "Ce post n'existe pas")
		return
	}

	// on charge les commentaires liés à ce post
	commentaires, err := app.db.GetCommentsByPostID(idPost, idUtilisateur)
	if err != nil {
		app.afficherErreur(w, r, 500, "Erreur serveur")
		return
	}

	categories := app.toutesLesCategories()

	donnees := &models.TemplateData{
		CurrentUser: utilisateur,
		Post:        post,
		Comments:    commentaires,
		Categories:  categories,
	}

	app.afficherPage(w, "post", donnees)
}
