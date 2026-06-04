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

	dsnURI2 := "db/threads.db"
	db2, err := sql.Open("sqlite", dsnURI2)
	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return
	}

	defer db2.Close()

	err = CreateThread(idUtilisateur, nomDuPoste, contenuDuTexte, nomDuLabel, db2)
	if (err != nil){
		fmt.Println(err)
	}
	nouvelIdFilDeDiscution := NombreElementDB(db2, "Threads")
	CreatePost(idUtilisateur, nouvelIdFilDeDiscution, contenuDuTexte, db, 0)

	RevenirSurLaPageAccueil(w, r, 0, false, true, nouvelIdFilDeDiscution)
}
