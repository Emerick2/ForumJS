package forumjs

import (
	"database/sql"
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

	listeLabel := ListeLabel()
	récupéréLesFilsDeDiscution := RécupéréLesFilsDeDiscution("Livre")
	fmt.Println(len(listeLabel))
	fmt.Println(len(récupéréLesFilsDeDiscution))
}

func ComplétéLaPageForum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("pages/discution.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
		return
	}

	données := map[string]interface{}{}

	err = tmpl.Execute(w, données)
	if err != nil {
		if isBrokenPipe(err) {
			return
		}
		fmt.Println("Erreur lors de l'exécution du template :", err)
	}
}

func ListeLabel() []string {
	listeLabel := []string{}

	dsnURI := "db/threads.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return listeLabel
	}
	defer db.Close()

	query := `
	SELECT DISTINCT label_name 
	FROM Threads
	ORDER BY label_name ASC`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Erreur :", err)
		return listeLabel
	}
	defer rows.Close()

	for rows.Next() {
		var unLabel string
		err := rows.Scan(
			&unLabel,
		)
		if err != nil {
			fmt.Println("Erreur :", err)
			return listeLabel
		}
		listeLabel = append(listeLabel, unLabel)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("Erreur :", err)
		return listeLabel
	}

	return listeLabel
}

func RécupéréLesFilsDeDiscution(recherche string) []Thread {
	listeThread := []Thread{}

	dsnURI := "db/threads.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return listeThread
	}
	defer db.Close()

	query := `
	SELECT id, user_id, name, message_content, label_name
	FROM Threads
	WHERE Label_name = ?
	ORDER BY id ASC`

	rows, err := db.Query(query, recherche)
	if err != nil {
		fmt.Println("Erreur :", err)
		return listeThread
	}
	defer rows.Close()

	for rows.Next() {
		var unThread Thread
		err := rows.Scan(
			&unThread.Id,
			&unThread.User_id,
			&unThread.Name,
			&unThread.Message_content,
			&unThread.Label_name,
		)
		if err != nil {
			fmt.Println("Erreur :", err)
			return listeThread
		}
		listeThread = append(listeThread, unThread)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("Erreur :", err)
		return listeThread
	}

	return listeThread
}
