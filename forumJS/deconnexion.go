package forumjs

import "net/http"

func Deconnecter(w http.ResponseWriter, r *http.Request) {
	SupprimerCookie(w)
	http.Redirect(w, r, "/Accueil", http.StatusSeeOther)
}
