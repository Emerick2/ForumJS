package handlers

import (
	"net/http"
	"strings"
)

// handleAjouterCommentaire ajoute un commentaire à un post
func (app *App) handleAjouterCommentaire(w http.ResponseWriter, r *http.Request) {
	// on récupère l'id du post depuis l'URL
	idPost := r.PathValue("id")

	// on récupère l'utilisateur connecté
	utilisateur := app.utilisateurConnecte(r)

	// on vérifie que le post existe dans la base de données
	post, err := app.db.GetPostByID(idPost, utilisateur.ID)
	if err != nil {
		app.afficherErreur(w, r, 500, "Erreur du serveur")
		return
	}
	if post == nil {
		app.afficherErreur(w, r, 404, "Ce post n'existe pas")
		return
	}

	// on lit le formulaire envoyé par l'utilisateur
	r.ParseForm()

	// on récupère le texte du commentaire
	texteCommentaire := strings.TrimSpace(r.FormValue("content"))

	// si le commentaire est vide on redirige sans rien faire
	if texteCommentaire == "" {
		http.Redirect(w, r, "/post/"+idPost, http.StatusSeeOther)
		return
	}

	// on enregistre le commentaire dans la base de données
	_, err = app.db.CreateComment(idPost, utilisateur.ID, texteCommentaire)
	if err != nil {
		app.afficherErreur(w, r, 500, "Impossible d'ajouter le commentaire")
		return
	}

	// on redirige vers le post avec l'ancre #comments pour voir le nouveau commentaire
	http.Redirect(w, r, "/post/"+idPost+"#comments", http.StatusSeeOther)
}

// handleSupprimerCommentaire supprime un commentaire
func (app *App) handleSupprimerCommentaire(w http.ResponseWriter, r *http.Request) {
	// on récupère l'id du commentaire depuis l'URL
	idCommentaire := r.PathValue("id")

	// on récupère l'utilisateur connecté
	utilisateur := app.utilisateurConnecte(r)

	// on récupère le commentaire depuis la base de données
	commentaire, err := app.db.GetCommentByID(idCommentaire, utilisateur.ID)
	if err != nil {
		app.afficherErreur(w, r, 500, "Erreur du serveur")
		return
	}
	if commentaire == nil {
		app.afficherErreur(w, r, 404, "Ce commentaire n'existe pas")
		return
	}

	// on vérifie que le commentaire appartient bien à l'utilisateur connecté
	if commentaire.UserID != utilisateur.ID {
		app.afficherErreur(w, r, 403, "Vous ne pouvez pas supprimer le commentaire de quelqu'un d'autre")
		return
	}

	// on supprime le commentaire de la base de données
	err = app.db.DeleteComment(idCommentaire)
	if err != nil {
		app.afficherErreur(w, r, 500, "Impossible de supprimer le commentaire")
		return
	}

	// on retourne sur le post
	http.Redirect(w, r, "/post/"+commentaire.PostID+"#comments", http.StatusSeeOther)
}
