package db

import (
	"database/sql"

	"forumJS/internal/models"
)

// CreateComment crée un commentaire sur un post
func (d *DB) CreateComment(idPost, idUtilisateur, contenu string) (*models.Comment, error) {
	id := nouvelID()

	_, err := d.conn.Exec(
		`INSERT INTO comments (id, post_id, user_id, content) VALUES (?, ?, ?, ?)`,
		id, idPost, idUtilisateur, contenu,
	)
	if err != nil {
		return nil, err
	}

	return d.GetCommentByID(id, "")
}

// GetCommentByID retourne un commentaire avec ses likes/dislikes
func (d *DB) GetCommentByID(id, idUtilisateurConnecte string) (*models.Comment, error) {
	c := &models.Comment{}

	err := d.conn.QueryRow(`
		SELECT
			cm.id, cm.post_id, cm.user_id, u.username, cm.content, cm.created_at, cm.updated_at,
			COALESCE(SUM(CASE WHEN r.value = 1  THEN 1 ELSE 0 END), 0) AS likes,
			COALESCE(SUM(CASE WHEN r.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes
		FROM comments cm
		JOIN users u ON u.id = cm.user_id
		LEFT JOIN reactions r ON r.target_id = cm.id AND r.target_type = 'comment'
		WHERE cm.id = ?
		GROUP BY cm.id`,
		id,
	).Scan(&c.ID, &c.PostID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.Likes, &c.Dislikes)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if idUtilisateurConnecte != "" {
		d.conn.QueryRow(
			`SELECT value FROM reactions WHERE user_id = ? AND target_id = ? AND target_type = 'comment'`,
			idUtilisateurConnecte, id,
		).Scan(&c.UserReaction)
	}

	return c, nil
}

// GetCommentsByPostID retourne tous les commentaires d'un post
func (d *DB) GetCommentsByPostID(idPost, idUtilisateurConnecte string) ([]models.Comment, error) {
	rows, err := d.conn.Query(`
		SELECT
			cm.id, cm.post_id, cm.user_id, u.username, cm.content, cm.created_at, cm.updated_at,
			COALESCE(SUM(CASE WHEN r.value = 1  THEN 1 ELSE 0 END), 0) AS likes,
			COALESCE(SUM(CASE WHEN r.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes
		FROM comments cm
		JOIN users u ON u.id = cm.user_id
		LEFT JOIN reactions r ON r.target_id = cm.id AND r.target_type = 'comment'
		WHERE cm.post_id = ?
		GROUP BY cm.id
		ORDER BY cm.created_at ASC`,
		idPost,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commentaires []models.Comment
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(
			&c.ID, &c.PostID, &c.UserID, &c.Username, &c.Content,
			&c.CreatedAt, &c.UpdatedAt, &c.Likes, &c.Dislikes,
		)
		if err != nil {
			return nil, err
		}

		if idUtilisateurConnecte != "" {
			d.conn.QueryRow(
				`SELECT value FROM reactions WHERE user_id = ? AND target_id = ? AND target_type = 'comment'`,
				idUtilisateurConnecte, c.ID,
			).Scan(&c.UserReaction)
		}

		commentaires = append(commentaires, c)
	}
	return commentaires, nil
}

// DeleteComment supprime un commentaire
func (d *DB) DeleteComment(id string) error {
	_, err := d.conn.Exec(`DELETE FROM comments WHERE id = ?`, id)
	return err
}
