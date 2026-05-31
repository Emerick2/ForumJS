package forumjs

import (
	"database/sql"
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
}

func CreatePost(userID int, threadID int, content string, db *sql.DB) error {
	requete := `
	INSERT INTO Posts (user_id, thread_id, content)
	VALUES (?, ?, ?)`

	_, err := db.Exec(requete, userID, threadID, content)
	return err
}

func GetPostsByThread(threadID int, db *sql.DB) ([]Post, error) {
	query := `
	SELECT id, user_id, thread_id, content, created_at, likes, dislikes 
	FROM Posts 
	WHERE thread_id = ?`

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