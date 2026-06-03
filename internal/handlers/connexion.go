package handlers

import (
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"forumJS/internal/models"
)

// durée de vie d'une session : 24 heures
const dureeSession = 24 * time.Hour

// nomCookie est le nom du cookie stocké dans le navigateur
const nomCookie = "session_id"

// chercherUtilisateurSession lit le cookie et retourne l'utilisateur connecté
// retourne nil si personne n'est connecté
func (app *App) chercherUtilisateurSession(r *http.Request) *models.User {
	cookie, err := r.Cookie(nomCookie)
	if err != nil {
		return nil
	}

	session, err := app.db.GetSession(cookie.Value)
	if err != nil || session == nil {
		return nil
	}

	// session expirée → on la supprime
	if time.Now().After(session.ExpiresAt) {
		app.db.DeleteSession(session.ID)
		return nil
	}

	utilisateur, _ := app.db.GetUserByID(session.UserID)
	return utilisateur
}

// creerCookieSession envoie le cookie de session au navigateur
func creerCookieSession(w http.ResponseWriter, session *models.Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     nomCookie,
		Value:    session.ID,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
}

// redigerSiDejaConnecte redirige vers "/" si l'utilisateur est déjà connecté
// retourne true si la redirection a eu lieu (le handler doit faire return)
func (app *App) redigerSiDejaConnecte(w http.ResponseWriter, r *http.Request) bool {
	if app.chercherUtilisateurSession(r) != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return true
	}
	return false
}

// handleFormulaireConnexion affiche la page de connexion
func (app *App) handleFormulaireConnexion(w http.ResponseWriter, r *http.Request) {
	if app.redigerSiDejaConnecte(w, r) {
		return
	}

	app.afficherPage(w, "login", &models.TemplateData{Categories: app.toutesLesCategories()})
}

// handleConnexion traite le formulaire de connexion (POST /login)
func (app *App) handleConnexion(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	email := strings.TrimSpace(r.FormValue("email"))
	motDePasse := r.FormValue("password")

	if email == "" || motDePasse == "" {
		app.afficherPage(w, "login", &models.TemplateData{
			Error:      "Email et mot de passe requis",
			Categories: app.toutesLesCategories(),
		})
		return
	}

	utilisateur, err := app.db.GetUserByEmail(email)
	if err != nil || utilisateur == nil {
		app.afficherPage(w, "login", &models.TemplateData{
			Error:      "Email ou mot de passe incorrect",
			Categories: app.toutesLesCategories(),
		})
		return
	}

	// on compare le mot de passe avec le hash stocké
	err = bcrypt.CompareHashAndPassword([]byte(utilisateur.PasswordHash), []byte(motDePasse))
	if err != nil {
		app.afficherPage(w, "login", &models.TemplateData{
			Error:      "Email ou mot de passe incorrect",
			Categories: app.toutesLesCategories(),
		})
		return
	}

	session, err := app.db.CreateSession(utilisateur.ID, dureeSession)
	if err != nil {
		app.afficherErreur(w, r, 500, "Erreur lors de la connexion")
		return
	}

	creerCookieSession(w, session)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// handleDeconnexion supprime la session et efface le cookie
func (app *App) handleDeconnexion(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(nomCookie)
	if err == nil {
		app.db.DeleteSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    nomCookie,
		Value:   "",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
		Path:    "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
