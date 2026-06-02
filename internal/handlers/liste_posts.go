package handlers

import (
	"net/http"

	"forum/internal/db"
	"forum/internal/models"
)

// handleIndex affiche la page d'accueil avec la liste des posts
// ticket #2 — Page de liste des postes
func (app *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	// si l'URL n'est pas exactement "/" on affiche une erreur 404
	if r.URL.Path != "/" {
		app.afficherErreur(w, r, 404, "Page introuvable")
		return
	}

	// on récupère l'utilisateur connecté (nil si pas connecté)
	utilisateur := app.chercherUtilisateurSession(r)

	// on récupère toutes les catégories pour la sidebar
	categories, _ := app.db.GetAllCategories()

	// on récupère les filtres depuis l'URL
	// exemples : /?category=sport  ou  /?filter=mine
	categorie := r.URL.Query().Get("category")
	filtre := r.URL.Query().Get("filter")

	// on prépare les critères de filtre
	filtrePosts := db.PostFilter{}

	if categorie != "" {
		filtrePosts.CategorySlug = categorie
	}

	// les filtres "mine" et "liked" nécessitent d'être connecté
	if utilisateur != nil {
		if filtre == "mine" {
			filtrePosts.UserID = utilisateur.ID
		}
		if filtre == "liked" {
			filtrePosts.LikedByUser = utilisateur.ID
		}
	}

	// ID de l'utilisateur connecté (vide si personne n'est connecté)
	idUtilisateur := ""
	if utilisateur != nil {
		idUtilisateur = utilisateur.ID
	}

	// on charge la liste des posts depuis la base de données
	listePosts, err := app.db.ListPosts(filtrePosts, idUtilisateur)
	if err != nil {
		app.afficherErreur(w, r, 500, "Impossible de charger les posts")
		return
	}

	donnees := &models.TemplateData{
		CurrentUser: utilisateur,
		Posts:       listePosts,
		Categories:  categories,
		SelectedCat: categorie,
		Filter:      filtre,
	}

	app.afficherPage(w, "index", donnees)
}
