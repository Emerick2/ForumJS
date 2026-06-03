package handlers

import (
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"forumJS/internal/models"
)

// handleFormulaireInscription affiche la page d'inscription
func (app *App) handleFormulaireInscription(w http.ResponseWriter, r *http.Request) {
	if app.redigerSiDejaConnecte(w, r) {
		return
	}

	app.afficherPage(w, "register", &models.TemplateData{Categories: app.toutesLesCategories()})
}

// handleInscription traite le formulaire d'inscription (POST /register)
func (app *App) handleInscription(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	email := strings.TrimSpace(strings.ToLower(r.FormValue("email")))
	pseudo := strings.TrimSpace(r.FormValue("username"))
	motDePasse := r.FormValue("password")
	confirmation := r.FormValue("confirm")

	// on valide les champs un par un
	if email == "" || pseudo == "" || motDePasse == "" {
		app.afficherPage(w, "register", &models.TemplateData{Error: "Tous les champs sont requis", Categories: app.toutesLesCategories()})
		return
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		app.afficherPage(w, "register", &models.TemplateData{Error: "L'email n'est pas valide", Categories: app.toutesLesCategories()})
		return
	}
	if len(pseudo) < 3 || len(pseudo) > 30 {
		app.afficherPage(w, "register", &models.TemplateData{Error: "Le pseudo doit faire entre 3 et 30 caractères", Categories: app.toutesLesCategories()})
		return
	}
	if len(motDePasse) < 6 {
		app.afficherPage(w, "register", &models.TemplateData{Error: "Le mot de passe doit faire au moins 6 caractères", Categories: app.toutesLesCategories()})
		return
	}
	if motDePasse != confirmation {
		app.afficherPage(w, "register", &models.TemplateData{Error: "Les deux mots de passe ne sont pas identiques", Categories: app.toutesLesCategories()})
		return
	}

	// on vérifie que l'email et le pseudo sont disponibles
	emailPris, _ := app.db.EmailExists(email)
	if emailPris {
		app.afficherPage(w, "register", &models.TemplateData{Error: "Cet email est déjà utilisé", Categories: app.toutesLesCategories()})
		return
	}

	pseudoPris, _ := app.db.UsernameExists(pseudo)
	if pseudoPris {
		app.afficherPage(w, "register", &models.TemplateData{Error: "Ce pseudo est déjà pris", Categories: app.toutesLesCategories()})
		return
	}

	// on hash le mot de passe avant de le stocker (jamais en clair)
	hashMotDePasse, err := bcrypt.GenerateFromPassword([]byte(motDePasse), 12)
	if err != nil {
		app.afficherErreur(w, r, 500, "Erreur serveur")
		return
	}

	nouvelUtilisateur, err := app.db.CreateUser(email, pseudo, string(hashMotDePasse))
	if err != nil {
		app.afficherErreur(w, r, 500, "Impossible de créer le compte")
		return
	}

	// on connecte directement l'utilisateur après l'inscription
	session, err := app.db.CreateSession(nouvelUtilisateur.ID, dureeSession)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	creerCookieSession(w, session)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
