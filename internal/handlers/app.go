package handlers

import (
	"html/template"
	"net/http"

	"forum/internal/db"
	"forum/internal/models"
)

// App contient la base de données et les templates HTML du site
type App struct {
	db        *db.DB
	templates map[string]*template.Template
}

// New crée l'application et charge tous les templates HTML
func New(baseDeDonnees *db.DB) (*App, error) {
	listeTemplates := make(map[string]*template.Template)

	// liste de toutes les pages du site
	// seuls les templates de nos tickets — les collègues ajouteront les leurs
	pages := []string{"index", "post", "login", "register", "error"}

	// pour chaque page on charge base.html + le fichier HTML de la page
	for _, nomPage := range pages {
		cheminBase := "web/templates/base.html"
		cheminPage := "web/templates/" + nomPage + ".html"

		fonctions := template.FuncMap{
			"add": func(a, b int) int { return a + b },
			"sub": func(a, b int) int { return a - b },
		}

		t, err := template.New("").Funcs(fonctions).ParseFiles(cheminBase, cheminPage)
		if err != nil {
			return nil, err
		}

		listeTemplates[nomPage] = t
	}

	monApp := &App{
		db:        baseDeDonnees,
		templates: listeTemplates,
	}
	return monApp, nil
}

// Routes enregistre toutes les URLs du site
func (app *App) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	// fichiers statiques : CSS, JavaScript, images
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// pages visibles par tout le monde (connecté ou non)
	mux.HandleFunc("GET /", app.handleIndex)
	mux.HandleFunc("GET /post/{id}", app.handleVoirPost)
	mux.HandleFunc("GET /login", app.handleFormulaireConnexion)
	mux.HandleFunc("POST /login", app.handleConnexion)
	mux.HandleFunc("GET /register", app.handleFormulaireInscription)
	mux.HandleFunc("POST /register", app.handleInscription)
	mux.HandleFunc("POST /logout", app.handleDeconnexion)

	// routes des collègues — à ajouter quand ils feront leurs PR
	// mux.HandleFunc("GET /post/new", ...)
	// mux.HandleFunc("POST /post/{id}/comment", ...)
	// mux.HandleFunc("POST /api/react", ...)

	return mux
}

// afficherPage envoie un template HTML au navigateur
func (app *App) afficherPage(w http.ResponseWriter, nomPage string, donnees *models.TemplateData) {
	t, existe := app.templates[nomPage]
	if !existe {
		http.Error(w, "page introuvable", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := t.ExecuteTemplate(w, "base", donnees)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// afficherErreur affiche la page d'erreur avec un code HTTP et un message
func (app *App) afficherErreur(w http.ResponseWriter, r *http.Request, codeErreur int, message string) {
	w.WriteHeader(codeErreur)

	donnees := &models.TemplateData{
		CurrentUser: app.chercherUtilisateurSession(r),
		Error:       message,
		ErrCode:     codeErreur,
	}

	app.afficherPage(w, "error", donnees)
}

// seulementConnecte redirige vers /login si l'utilisateur n'est pas connecté
func (app *App) seulementConnecte(suite http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utilisateur := app.chercherUtilisateurSession(r)

		if utilisateur == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		suite(w, r)
	}
}
