package forumjs

import (
	"database/sql"
	"fmt"
	"time"
)

type Post struct {
	Id        int
	UserId    int
	ThreadId  int
	Content   string
	CreatedAt time.Time
	Likes     int
	Dislikes  int
	Answer    int
}

type PostTableauDeBord struct {
	Id            int
	UserId        int
	ThreadId      int
	Content       string
	CreatedAt     time.Time
	CreatedAtText string
	Likes         int
	Dislikes      int
	Answer        int
	NameUser      string
	IconeLike     string
	IconeDislike  string
}

type Thread struct {
	Id              int
	Name            string
	User_id         int
	Message_content string
	Label_name      string
}

type User struct {
	Id            int
	CreatedAt     time.Time
	CreatedAtText string
	Name          string
	Adresse_email string
	Password      string
}

func CreatePost(userID int, threadID int, content string, db *sql.DB, answer int) error {
	requete := `
	INSERT INTO Posts (user_id, thread_id, content, answer)
	VALUES (?, ?, ?, ?)`

	_, err := db.Exec(requete, userID, threadID, content, answer)
	return err
}

func CreateThread(idUtilisateur int, nomDuLabel string, contenuDuTexte string, label_name string, db *sql.DB) error {
	dsnURI := "db/threads.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return err
	}
	defer db.Close()

	requete := `
	INSERT INTO Threads (user_id, name, message_content, label_name)
	VALUES (?, ?, ?, ?)`

	_, err = db.Exec(requete, idUtilisateur, nomDuLabel, contenuDuTexte, label_name)
	return err
}

func GetPostsByThread(threadID int, db *sql.DB) ([]Post, error) {
	query := `
	SELECT id, user_id, thread_id, content, created_at, likes, dislikes, answer 
	FROM Posts 
	WHERE thread_id = ? ORDER BY created_at ASC`

	rows, err := db.Query(query, threadID)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		listePosts = append(listePosts, unPost)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return listePosts, nil
}

func GetThread() ([]Thread, error) {
	dsnURI := "db/threads.db"
	db, err := sql.Open("sqlite", dsnURI)
	if err != nil {
		fmt.Println("Erreur d'ouverture :", err)
		return nil, err
	}
	defer db.Close()

	query := `
	SELECT id, user_id, name, message_content, label_name 
	FROM Threads
	ORDER BY id ASC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	listeThread := []Thread{}

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
			return nil, err
		}
		listeThread = append(listeThread, unThread)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return listeThread, nil
}

func GetPost() []Post {
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
	ORDER BY created_at DESC;`

	rows, err := db.Query(query)
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

	return listePosts
}

func NombreElementDB(db *sql.DB, nomTable string) int {
	query := `
	SELECT id 
	FROM Threads `

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	defer rows.Close()

	total := 0
	for rows.Next() {
		total++
	}
	fmt.Println(total)
	return total
}
