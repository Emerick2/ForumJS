package forumjs

import (
	"fmt"
	"net/http"
	"text/template"
)

func ComplétéLaPageAccueil(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	idUtilisateur := VérifierCookie(r)
	nomAAfficher := "Invité"

	if idUtilisateur != 0 {
		utilisateur := VoirUtilisateurs(idUtilisateur)
		if utilisateur.nom != "" {
			nomAAfficher = utilisateur.nom
		}
	}

	données := map[string]interface{}{
		"NomUtilisateur": nomAAfficher,
	}

	tmpl, err := template.ParseFiles("pages/main.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, données)
	if err != nil {
		if isBrokenPipe(err) {
			return
		}
		fmt.Println("Erreur lors de l'exécution du template :", err)
	}
}
