package forumjs

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"
)

// func TableauDeBord(w http.ResponseWriter, r *http.Request) {
// 	nombreTotalAimeSurCommentaires := NombreTotalAimeSurCommentaires()
// 	nombreTotalMessagePublier := NombreTotalMessagePublier()
// 	nombreTotalUtilisateur := NombreTotalUtilisateur()
// 	derniersMessagesPublié := DerniersMessagesPublié(5)
// 	derniersUtilisateursCréé := DerniersUtilisateursCréé(5)

// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	tmpl, err := template.ParseFiles("pages/tableau-de-bord.html")
// 	if err != nil {
// 		http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
// 		return
// 	}

// 	données := map[string]interface{}{
// 		"NombreTotalAimeSurCommentaires": nombreTotalAimeSurCommentaires,
// 		"NombreTotalMessagePublier":      nombreTotalMessagePublier,
// 		"NombreTotalUtilisateur":         nombreTotalUtilisateur,
// 	}

// 	err = tmpl.Execute(w, données)
// 	if err != nil {
// 		if isBrokenPipe(err) {
// 			return
// 		}
// 		fmt.Println("Erreur lors de l'exécution du template :", err)
// 	}

// 	for i := 0; i < len(derniersMessagesPublié); i++ {
// 		AfficherPost(derniersMessagesPublié[i], w, r, false, 0, false)
// 	}

// 	for i := 0; i < len(derniersUtilisateursCréé); i++ {
// 		AfficherUtilisateur(derniersUtilisateursCréé[i], w, r)
// 	}

// 	/*
// 		*Les derniers messages publié.
// 		Le nombre total de j'aime mis sur les commentaires
// 		Le nombre total de message publier
// 		Le nombre total d'utilisateur
// 		Les fils de discution triés par ceux avec le plus de commentaires

// 		une listes pour voirs tous les utilisateurs du site.
// 	*/
// }

func TableauDeBord(w http.ResponseWriter, r *http.Request) {
	nombreTotalAimeSurCommentaires := NombreTotalAimeSurCommentaires()
	nombreTotalMessagePublier := NombreTotalMessagePublier()
	nombreTotalUtilisateur := NombreTotalUtilisateur()
	derniersMessagesPublié := DerniersMessagesPublié(5, w, r)
	derniersUtilisateursCréé := DerniersUtilisateursCréé(5)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("pages/tableau-de-bord.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
		return
	}

	données := map[string]interface{}{
		"NombreTotalAimeSurCommentaires": nombreTotalAimeSurCommentaires,
		"NombreTotalMessagePublier":      nombreTotalMessagePublier,
		"NombreTotalUtilisateur":         nombreTotalUtilisateur,
		"DerniersPosts":                  derniersMessagesPublié,
		"DerniersUsers":                  derniersUtilisateursCréé,
	}

	err = tmpl.Execute(w, données)
	if err != nil {
		if isBrokenPipe(err) {
			return
		}
		fmt.Println("Erreur lors de l'exécution du template :", err)
	}
}

func OuvrirDB(dsnURI string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return db, nil
}

func ConnaitreNombre(rows *sql.Rows) int {
	total := 0
	valeur := 0
	for rows.Next() {
		err := rows.Scan(&valeur)
		if err != nil {
			fmt.Println(err)
		}
		total += valeur
	}
	return total
}

func DerniersMessagesPublié(nombreMaximum int, w http.ResponseWriter, r *http.Request) []PostTableauDeBord {
	// db, err := OuvrirDB("db/forum.db")
	dsnURI := "db/forum.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		return nil
	}
	defer db.Close()

	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return nil
	}

	query := `
	SELECT id, user_id, thread_id, content, created_at, likes, dislikes, answer 
	FROM Posts 
	ORDER BY created_at DESC
	LIMIT ?;`

	rows, err := db.Query(query, nombreMaximum)
	if err != nil {
		fmt.Println("Erreur :", err)
		return nil
	}
	defer rows.Close()

	listePosts := []PostTableauDeBord{}

	for rows.Next() {
		var unPost PostTableauDeBord
		err := rows.Scan(
			&unPost.Id,
			&unPost.UserId,
			&unPost.ThreadId,
			&unPost.Content,
			&unPost.CreatedAt,
			&unPost.Likes,
			&unPost.Dislikes,
			&unPost.Answer,
		)
		if err != nil {
			fmt.Println("Erreur :", err)
			return nil
		}
		unPost.NameUser = "Compte suprimé"
		valeur := VoirUtilisateurs(unPost.UserId)
		if valeur.nom != "" {
			unPost.NameUser = valeur.nom
		}

		unPost.CreatedAtText = Date(unPost.CreatedAt)

		unPost.IconeLike = "/images/aime.svg"
		unPost.IconeDislike = "/images/aime.svg"

		idUtilisateur := VérifierCookie(r)
		if idUtilisateur != 0 {
			if LireTableauInteractionUtilisateur(w, r, idUtilisateur, unPost.Id, unPost.ThreadId, "likes") {
				unPost.IconeLike = "/images/aimeActif.svg"
			}
			if LireTableauInteractionUtilisateur(w, r, idUtilisateur, unPost.Id, unPost.ThreadId, "dislikes") {
				unPost.IconeDislike = "/images/aimeActif.svg"
			}
		}

		listePosts = append(listePosts, unPost)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("Erreur :", err)
		return nil
	}

	fmt.Println(len(listePosts))
	return listePosts
}

func NombreTotalAimeSurCommentaires() int {
	// db, err := OuvrirDB("db/forum.db")
	dsnURI := "db/forum.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		return 0
	}
	defer db.Close()

	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return 0
	}

	query := `
	SELECT SUM(likes)
	FROM Posts;`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Erreur :", err)
		return 0
	}
	defer rows.Close()

	return ConnaitreNombre(rows)
}

func NombreTotalMessagePublier() int {
	// db, err := OuvrirDB("db/forum.db")
	dsnURI := "db/forum.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		return 0
	}
	defer db.Close()

	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return 0
	}

	query := `
	SELECT COUNT(id)
	FROM Posts;`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Erreur :", err)
		return 0
	}
	defer rows.Close()

	return ConnaitreNombre(rows)
}

func NombreTotalUtilisateur() int {
	dsnURI := "db/user.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		return 0
	}
	defer db.Close()

	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return 0
	}

	query := `
	SELECT COUNT(UserId)
	FROM user;`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Erreur :", err)
		return 0
	}
	defer rows.Close()

	return ConnaitreNombre(rows)
}

func DerniersUtilisateursCréé(limite int) []User {
	// db, err := OuvrirDB("db/forum.db")
	dsnURI := "db/user.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		return nil
	}
	defer db.Close()

	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return nil
	}

	query := `
	SELECT UserId, Email, NomUtilisateur, CreatedAt 
	FROM user 
	ORDER BY CreatedAt DESC
	LIMIT ?;`

	rows, err := db.Query(query, limite)
	if err != nil {
		fmt.Println("Erreur :", err)
		return nil
	}
	defer rows.Close()

	listePosts := []User{}

	for rows.Next() {
		var unPost User
		err := rows.Scan(
			&unPost.Id,
			&unPost.Adresse_email,
			&unPost.Name,
			&unPost.CreatedAt,
		)
		if err != nil {
			fmt.Println("Erreur :", err)
			return nil
		}
		unPost.CreatedAtText = Date(unPost.CreatedAt)
		listePosts = append(listePosts, unPost)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("Erreur :", err)
		return nil
	}

	// fmt.Println(len(listePosts))
	return listePosts
}
