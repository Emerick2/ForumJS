package forumjs

import (
	"database/sql"
	"fmt"
	"net/http"
)

func NouveauFilDeDiscution(w http.ResponseWriter, r *http.Request) {
	idUtilisateur := VérifierCookie(r)
	if idUtilisateur == 0 {
		return
	}

	nomDuPoste := r.FormValue("nomDuPoste")
	nomDuLabel := r.FormValue("nomDuLabel")
	contenuDuTexte := r.FormValue("contenuDuTexte")

	dsnURI := "db/forum.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return
	}

	defer db.Close()

	CreateThread(idUtilisateur, nomDuPoste, contenuDuTexte, nomDuLabel, db)
	nouvelIdFilDeDiscution := NombreElementDB(db, "Threads")
	CreatePost(idUtilisateur, nouvelIdFilDeDiscution, contenuDuTexte, db, 0)

	RevenirSurLaPageAccueil(w, r, 0, false, true, nouvelIdFilDeDiscution)
}
