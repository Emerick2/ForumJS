package forumjs

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"
)

func TableauDeBord(w http.ResponseWriter, r *http.Request) {
	DerniersMessagesPublié(5)
	nombreTotalAimeSurCommentaires := NombreTotalAimeSurCommentaires()
	nombreTotalMessagePublier := NombreTotalMessagePublier()
	nombreTotalUtilisateur := NombreTotalUtilisateur()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("pages/tableau-de-bord.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
		return
	}

	données := map[string]interface{}{
		"NombreTotalAimeSurCommentaires": nombreTotalAimeSurCommentaires,
		"NombreTotalMessagePublier": nombreTotalMessagePublier,
		"NombreTotalUtilisateur": nombreTotalUtilisateur,
	}

	err = tmpl.Execute(w, données)
	if err != nil {
		if isBrokenPipe(err) {
			return
		}
		fmt.Println("Erreur lors de l'exécution du template :", err)
	}

	/*
		*Les derniers messages publié.
		Le nombre total de j'aime mis sur les commentaires
		Le nombre total de message publier
		Le nombre total d'utilisateur
		Les fils de discution triés par ceux avec le plus de commentaires

		une listes pour voirs tous les utilisateurs du site.
	*/
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

func DerniersMessagesPublié(nombreMaximum int) []Post {
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

	listePosts := []Post{}

	for rows.Next() {
		var unPost Post
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
	SELECT COUNT(likes)
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
