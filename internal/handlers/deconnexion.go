func handleDeconnexion(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie {
		Name : "session_utilisateur",
		MaxAge : -1,
		Path : "/",
	}
	http.SetCookie(w, cookie) 
	http.Redirect(w, r, "/", http.StatusSeeOther)
}