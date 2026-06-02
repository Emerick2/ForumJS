package db

import "forum/internal/models"

// PostFilter contient les critères de filtre pour la liste des posts
type PostFilter struct {
	CategorySlug string // filtre par catégorie
	UserID       string // filtre : posts créés par cet utilisateur
	LikedByUser  string // filtre : posts aimés par cet utilisateur
}

// ListPosts retourne la liste des posts avec les filtres appliqués
func (d *DB) ListPosts(filtre PostFilter, idUtilisateurConnecte string) ([]models.Post, error) {
	requete := `
		SELECT
			p.id, p.user_id, u.username, p.title, p.content, p.created_at, p.updated_at,
			COALESCE(SUM(CASE WHEN r.value = 1  THEN 1 ELSE 0 END), 0) AS likes,
			COALESCE(SUM(CASE WHEN r.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes,
			(SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) AS nb_commentaires
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN reactions r ON r.target_id = p.id AND r.target_type = 'post'`

	var arguments []interface{}
	var conditions []string

	// filtre par catégorie
	if filtre.CategorySlug != "" {
		requete += ` JOIN post_categories pc ON pc.post_id = p.id
		             JOIN categories cat ON cat.id = pc.category_id`
		conditions = append(conditions, "cat.slug = ?")
		arguments = append(arguments, filtre.CategorySlug)
	}

	// filtre "mes posts"
	if filtre.UserID != "" {
		conditions = append(conditions, "p.user_id = ?")
		arguments = append(arguments, filtre.UserID)
	}

	// filtre "posts aimés"
	if filtre.LikedByUser != "" {
		requete += ` JOIN reactions rl ON rl.target_id = p.id AND rl.target_type = 'post' AND rl.value = 1`
		conditions = append(conditions, "rl.user_id = ?")
		arguments = append(arguments, filtre.LikedByUser)
	}

	if len(conditions) > 0 {
		requete += " WHERE "
		for i, condition := range conditions {
			if i > 0 {
				requete += " AND "
			}
			requete += condition
		}
	}

	requete += " GROUP BY p.id ORDER BY p.created_at DESC"

	rows, err := d.conn.Query(requete, arguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err := rows.Scan(
			&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content,
			&p.CreatedAt, &p.UpdatedAt, &p.Likes, &p.Dislikes, &p.CommentCount,
		)
		if err != nil {
			return nil, err
		}

		if idUtilisateurConnecte != "" {
			d.conn.QueryRow(
				`SELECT value FROM reactions WHERE user_id = ? AND target_id = ? AND target_type = 'post'`,
				idUtilisateurConnecte, p.ID,
			).Scan(&p.UserReaction)
		}

		cats, _ := d.getCategoriesPost(p.ID)
		p.Categories = cats

		posts = append(posts, p)
	}
	return posts, nil
}
