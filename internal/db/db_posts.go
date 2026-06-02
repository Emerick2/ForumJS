package db

import (
	"database/sql"

	"forumJS/internal/models"
)

// CreatePost crée un nouveau post avec ses catégories
// on utilise une transaction : si une étape échoue, tout est annulé
func (d *DB) CreatePost(idUtilisateur, titre, contenu string, idCategories []int) (*models.Post, error) {
	idPost := nouvelID()

	tx, err := d.conn.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(
		`INSERT INTO posts (id, user_id, title, content) VALUES (?, ?, ?, ?)`,
		idPost, idUtilisateur, titre, contenu,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, idCat := range idCategories {
		_, err = tx.Exec(
			`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`, idPost, idCat,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return d.GetPostByID(idPost, "")
}

// GetPostByID retourne un post avec ses likes/dislikes et le nombre de commentaires
func (d *DB) GetPostByID(idPost, idUtilisateurConnecte string) (*models.Post, error) {
	p := &models.Post{}

	err := d.conn.QueryRow(`
		SELECT
			p.id, p.user_id, u.username, p.title, p.content, p.created_at, p.updated_at,
			COALESCE(SUM(CASE WHEN r.value = 1  THEN 1 ELSE 0 END), 0) AS likes,
			COALESCE(SUM(CASE WHEN r.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes,
			(SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) AS nb_commentaires
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN reactions r ON r.target_id = p.id AND r.target_type = 'post'
		WHERE p.id = ?
		GROUP BY p.id`,
		idPost,
	).Scan(
		&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content,
		&p.CreatedAt, &p.UpdatedAt, &p.Likes, &p.Dislikes, &p.CommentCount,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// on regarde si l'utilisateur connecté a déjà voté sur ce post
	if idUtilisateurConnecte != "" {
		d.conn.QueryRow(
			`SELECT value FROM reactions WHERE user_id = ? AND target_id = ? AND target_type = 'post'`,
			idUtilisateurConnecte, idPost,
		).Scan(&p.UserReaction)
	}

	categories, err := d.getCategoriesPost(idPost)
	if err != nil {
		return nil, err
	}
	p.Categories = categories

	return p, nil
}

// getCategoriesPost retourne les catégories associées à un post
func (d *DB) getCategoriesPost(idPost string) ([]models.Category, error) {
	rows, err := d.conn.Query(`
		SELECT c.id, c.name, c.slug
		FROM categories c
		JOIN post_categories pc ON pc.category_id = c.id
		WHERE pc.post_id = ?`,
		idPost,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Slug); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

// UpdatePost modifie le titre, le contenu et les catégories d'un post
func (d *DB) UpdatePost(idPost, titre, contenu string, idCategories []int) error {
	tx, err := d.conn.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`UPDATE posts SET title = ?, content = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		titre, contenu, idPost,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`DELETE FROM post_categories WHERE post_id = ?`, idPost)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, idCat := range idCategories {
		_, err = tx.Exec(
			`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`, idPost, idCat,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// DeletePost supprime un post (commentaires et réactions supprimés automatiquement via CASCADE)
func (d *DB) DeletePost(idPost string) error {
	_, err := d.conn.Exec(`DELETE FROM posts WHERE id = ?`, idPost)
	return err
}
