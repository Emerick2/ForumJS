package main

type Post struct {
	Id        int
	UserId    int
	ThreadId  int
	Content   string
	CreatedAt string
	Likes     int
	Dislikes  int
}

func createPost(UserId int, ThreadId int, Content string) error {
	requete := `
	INSERT INTO Posts (user_id, thread_id, content)
	VALUES (?, ?, ?)`

	_, err := db.Exec(requete, UserId, ThreadId, Content)
	if err != nil {
		return err
	}
	return nil
}

func recupPost(ThreadId int) ([]Post, error) {
	rows, err := db.Query(`
	SELECT * FROM Posts WHERE thread_id = ?`, ThreadId)

	if err != nil {
		return nil, err
	}
	// defer rows.close()

	listePosts := []Post{}

	for rows.Next() {
		unPost := Post{}
		err := rows.Scan(&unPost.Id, &unPost.UserId, &unPost.ThreadId, &unPost.Content, &unPost.CreatedAt, &unPost.Likes, &unPost.Dislikes)
		if err != nil {
			return nil, err
		}
		listePosts = append(listePosts, unPost)
	}
	return listePosts, nil
}
