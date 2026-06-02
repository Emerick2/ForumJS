package db

import (
	"database/sql"
	"time"

	"forum/internal/models"
)

// ─── Utilisateurs ─────────────────────────────────────────────────────────────

// CreateUser crée un nouvel utilisateur
func (d *DB) CreateUser(email, username, hash string) (*models.User, error) {
	id := nouvelID()

	_, err := d.conn.Exec(
		`INSERT INTO users (id, email, username, password_hash) VALUES (?, ?, ?, ?)`,
		id, email, username, hash,
	)
	if err != nil {
		return nil, err
	}

	return d.GetUserByID(id)
}

// GetUserByEmail cherche un utilisateur par son email
func (d *DB) GetUserByEmail(email string) (*models.User, error) {
	u := &models.User{}

	err := d.conn.QueryRow(
		`SELECT id, email, username, password_hash, created_at FROM users WHERE email = ?`, email,
	).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

// GetUserByID cherche un utilisateur par son ID
func (d *DB) GetUserByID(id string) (*models.User, error) {
	u := &models.User{}

	err := d.conn.QueryRow(
		`SELECT id, email, username, password_hash, created_at FROM users WHERE id = ?`, id,
	).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

// EmailExists vérifie si un email est déjà utilisé
func (d *DB) EmailExists(email string) (bool, error) {
	var nombre int
	err := d.conn.QueryRow(`SELECT COUNT(*) FROM users WHERE email = ?`, email).Scan(&nombre)
	return nombre > 0, err
}

// UsernameExists vérifie si un pseudo est déjà pris
func (d *DB) UsernameExists(username string) (bool, error) {
	var nombre int
	err := d.conn.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ?`, username).Scan(&nombre)
	return nombre > 0, err
}

// ─── Sessions ─────────────────────────────────────────────────────────────────

// CreateSession crée une session et supprime les anciennes (une seule session par utilisateur)
func (d *DB) CreateSession(idUtilisateur string, duree time.Duration) (*models.Session, error) {
	_, err := d.conn.Exec(`DELETE FROM sessions WHERE user_id = ?`, idUtilisateur)
	if err != nil {
		return nil, err
	}

	session := &models.Session{
		ID:        nouvelID(),
		UserID:    idUtilisateur,
		ExpiresAt: time.Now().Add(duree),
	}

	_, err = d.conn.Exec(
		`INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)`,
		session.ID, session.UserID, session.ExpiresAt,
	)
	return session, err
}

// GetSession cherche une session par son ID
func (d *DB) GetSession(id string) (*models.Session, error) {
	s := &models.Session{}

	err := d.conn.QueryRow(
		`SELECT id, user_id, expires_at FROM sessions WHERE id = ?`, id,
	).Scan(&s.ID, &s.UserID, &s.ExpiresAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}

// DeleteSession supprime une session
func (d *DB) DeleteSession(id string) error {
	_, err := d.conn.Exec(`DELETE FROM sessions WHERE id = ?`, id)
	return err
}

// ─── Catégories ───────────────────────────────────────────────────────────────

// GetAllCategories retourne toutes les catégories
func (d *DB) GetAllCategories() ([]models.Category, error) {
	rows, err := d.conn.Query(`SELECT id, name, slug FROM categories ORDER BY name`)
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
