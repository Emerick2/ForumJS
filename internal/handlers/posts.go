package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"forum/internal/db"
	"forum/internal/models"
)

// handleIndex affiche la page d'accueil avec la liste des posts
func (app *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	// si l'URL n'est pas exactement "/" on affiche une erreur 404
	if r.URL.Path != "/" {
		app.afficherErreur(w, r, 404, "Page introuvable")
		return
	}

	// on récupère l'utilisateur connecté (peut être nil si pas connecté)
	utilisateur := app.utilisateurConnecte(r)

	// on récupère toutes les catégories pour la sidebar
	categories, _ := app.db.GetAllCategories()

	// on regarde les filtres dans l'URL
	// exemple : /?category=sport ou /?filter=mine
	categorie := r.URL.Query().Get("category")
	filtre := r.URL.Query().Get("filter")

	// on prépare le filtre pour la base de données
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

	// on récupère l'ID de l'utilisateur (vide si pas connecté)
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

	// on envoie tout ça au template HTML
	donnees := &models.TemplateData{
		CurrentUser: utilisateur,
		Posts:       listePosts,
		Categories:  categories,
		SelectedCat: categorie,
		Filter:      filtre,
	}

	app.afficherPage(w, "index", donnees)
}

// handleVoirPost affiche un post avec ses commentaires
func (app *App) handleVoirPost(w http.ResponseWriter, r *http.Request) {
	// on récupère l'id du post depuis l'URL
	idPost := r.PathValue("id")

	utilisateur := app.utilisateurConnecte(r)

	idUtilisateur := ""
	if utilisateur != nil {
		idUtilisateur = utilisateur.ID
	}

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

	// on charge les commentaires du post
	commentaires, err := app.db.GetCommentsByPostID(idPost, idUtilisateur)
	if err != nil {
		app.afficherErreur(w, r, 500, "Erreur serveur")
		return
	}

	categories, _ := app.db.GetAllCategories()

	donnees := &models.TemplateData{
		CurrentUser: utilisateur,
		Post:        post,
		Comments:    commentaires,
		Categories:  categories,
	}

	app.afficherPage(w, "post", donnees)
}

// handleFormulaireNouveauPost affiche le formulaire de création de post
func (app *App) handleFormulaireNouveauPost(w http.ResponseWriter, r *http.Request) {
	categories, _ := app.db.GetAllCategories()

	donnees := &models.TemplateData{
		CurrentUser: app.utilisateurConnecte(r),
		Categories:  categories,
	}

	app.afficherPage(w, "create_post", donnees)
}

// handleCreerPost traite le formulaire et enregistre le nouveau post
func (app *App) handleCreerPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	titre := strings.TrimSpace(r.FormValue("title"))
	contenu := strings.TrimSpace(r.FormValue("content"))

	// on récupère les catégories cochées dans le formulaire
	idCategories := recupererIdCategories(r.Form["categories"])

	utilisateur := app.utilisateurConnecte(r)
	categories, _ := app.db.GetAllCategories()

	// vérifications
	if titre == "" || contenu == "" {
		donnees := &models.TemplateData{
			CurrentUser: utilisateur,
			Categories:  categories,
			Error:       "Le titre et le contenu sont obligatoires",
		}
		app.afficherPage(w, "create_post", donnees)
		return
	}

	if len(titre) > 200 {
		donnees := &models.TemplateData{
			CurrentUser: utilisateur,
			Categories:  categories,
			Error:       "Le titre est trop long (200 caractères maximum)",
		}
		app.afficherPage(w, "create_post", donnees)
		return
	}

	if len(idCategories) == 0 {
		donnees := &models.TemplateData{
			CurrentUser: utilisateur,
			Categories:  categories,
			Error:       "Choisissez au moins une catégorie",
		}
		app.afficherPage(w, "create_post", donnees)
		return
	}

	// on sauvegarde le post dans la base de données
	nouveauPost, err := app.db.CreatePost(utilisateur.ID, titre, contenu, idCategories)
	if err != nil {
		app.afficherErreur(w, r, 500, "Impossible de créer le post")
		return
	}

	// on redirige vers le post créé
	http.Redirect(w, r, "/post/"+nouveauPost.ID, http.StatusSeeOther)
}

// handleFormulaireModifierPost affiche le formulaire de modification
func (app *App) handleFormulaireModifierPost(w http.ResponseWriter, r *http.Request) {
	idPost := r.PathValue("id")
	utilisateur := app.utilisateurConnecte(r)

	// on charge le post existant
	post, err := app.db.GetPostByID(idPost, utilisateur.ID)
	if err != nil || post == nil {
		app.afficherErreur(w, r, 404, "Post introuvable")
		return
	}

	// on vérifie que le post appartient à l'utilisateur
	if post.UserID != utilisateur.ID {
		app.afficherErreur(w, r, 403, "Vous ne pouvez modifier que vos propres posts")
		return
	}

	categories, _ := app.db.GetAllCategories()

	donnees := &models.TemplateData{
		CurrentUser: utilisateur,
		Post:        post,
		Categories:  categories,
	}

	app.afficherPage(w, "edit_post", donnees)
}

// handleModifierPost sauvegarde les modifications du post
func (app *App) handleModifierPost(w http.ResponseWriter, r *http.Request) {
	idPost := r.PathValue("id")
	utilisateur := app.utilisateurConnecte(r)

	// on vérifie que le post existe et appartient à l'utilisateur
	post, err := app.db.GetPostByID(idPost, utilisateur.ID)
	if err != nil || post == nil {
		app.afficherErreur(w, r, 404, "Post introuvable")
		return
	}

	if post.UserID != utilisateur.ID {
		app.afficherErreur(w, r, 403, "Vous ne pouvez modifier que vos propres posts")
		return
	}

	r.ParseForm()

	titre := strings.TrimSpace(r.FormValue("title"))
	contenu := strings.TrimSpace(r.FormValue("content"))
	idCategories := recupererIdCategories(r.Form["categories"])

	categories, _ := app.db.GetAllCategories()

	if titre == "" || contenu == "" {
		donnees := &models.TemplateData{
			CurrentUser: utilisateur,
			Post:        post,
			Categories:  categories,
			Error:       "Le titre et le contenu sont obligatoires",
		}
		app.afficherPage(w, "edit_post", donnees)
		return
	}

	if len(idCategories) == 0 {
		donnees := &models.TemplateData{
			CurrentUser: utilisateur,
			Post:        post,
			Categories:  categories,
			Error:       "Choisissez au moins une catégorie",
		}
		app.afficherPage(w, "edit_post", donnees)
		return
	}

	// on met à jour le post dans la base de données
	err = app.db.UpdatePost(idPost, titre, contenu, idCategories)
	if err != nil {
		app.afficherErreur(w, r, 500, "Impossible de modifier le post")
		return
	}

	http.Redirect(w, r, "/post/"+idPost, http.StatusSeeOther)
}

// handleSupprimerPost supprime un post
func (app *App) handleSupprimerPost(w http.ResponseWriter, r *http.Request) {
	idPost := r.PathValue("id")
	utilisateur := app.utilisateurConnecte(r)

	// on vérifie que le post existe et appartient à l'utilisateur
	post, err := app.db.GetPostByID(idPost, utilisateur.ID)
	if err != nil || post == nil {
		app.afficherErreur(w, r, 404, "Post introuvable")
		return
	}

	if post.UserID != utilisateur.ID {
		app.afficherErreur(w, r, 403, "Vous ne pouvez supprimer que vos propres posts")
		return
	}

	// on supprime le post
	err = app.db.DeletePost(idPost)
	if err != nil {
		app.afficherErreur(w, r, 500, "Impossible de supprimer le post")
		return
	}

	// on retourne à l'accueil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// recupererIdCategories convertit une liste de strings en liste d'entiers
func recupererIdCategories(valeurs []string) []int {
	var ids []int
	for _, v := range valeurs {
		nombre, err := strconv.Atoi(v)
		if err == nil {
			ids = append(ids, nombre)
		}
	}
	return ids
}
