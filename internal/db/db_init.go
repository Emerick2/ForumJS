package db

import (
	"crypto/rand"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // driver SQLite (pas besoin de GCC)
)

// DB est notre structure qui contient la connexion à la base de données
type DB struct {
	conn *sql.DB
}

// Init ouvre la base de données et crée les tables si elles n'existent pas
func Init(chemin string) (*DB, error) {
	// _foreign_keys=on active les liens entre tables (clés étrangères)
	connexion, err := sql.Open("sqlite", chemin+"?_foreign_keys=on")
	if err != nil {
		return nil, err
	}

	err = connexion.Ping()
	if err != nil {
		return nil, err
	}

	maDB := &DB{conn: connexion}

	err = maDB.creerTables()
	if err != nil {
		return nil, err
	}

	return maDB, nil
}

// Close ferme la connexion à la base de données
func (d *DB) Close() {
	d.conn.Close()
}

// creerTables crée toutes les tables au démarrage
func (d *DB) creerTables() error {
	requetes := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			expires_at DATETIME NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			slug TEXT UNIQUE NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS posts (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS post_categories (
			post_id TEXT NOT NULL,
			category_id INTEGER NOT NULL,
			PRIMARY KEY (post_id, category_id),
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS comments (
			id TEXT PRIMARY KEY,
			post_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		// PRIMARY KEY (user_id, target_id, target_type) empêche de voter deux fois
		`CREATE TABLE IF NOT EXISTS reactions (
			user_id TEXT NOT NULL,
			target_id TEXT NOT NULL,
			target_type TEXT NOT NULL CHECK(target_type IN ('post','comment')),
			value INTEGER NOT NULL CHECK(value IN (-1,1)),
			PRIMARY KEY (user_id, target_id, target_type),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
	}

	for _, requete := range requetes {
		_, err := d.conn.Exec(requete)
		if err != nil {
			return err
		}
	}

	return d.ajouterCategories()
}

// ajouterCategories insère les catégories par défaut
func (d *DB) ajouterCategories() error {
	categories := []struct{ nom, slug string }{
		{"Général", "general"},
		{"Technologie", "technologie"},
		{"Science", "science"},
		{"Jeux vidéo", "jeux-video"},
		{"Cinéma & Séries", "cinema-series"},
		{"Musique", "musique"},
		{"Sport", "sport"},
		{"Politique", "politique"},
		{"Humour", "humour"},
		{"Actualités", "actualites"},
	}

	for _, cat := range categories {
		// INSERT OR IGNORE = on n'insère pas si la catégorie existe déjà
		_, err := d.conn.Exec(
			`INSERT OR IGNORE INTO categories (name, slug) VALUES (?, ?)`,
			cat.nom, cat.slug,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// nouvelID génère un identifiant unique aléatoire
func nouvelID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// NewID est la version publique de nouvelID
func NewID() string {
	return nouvelID()
}
