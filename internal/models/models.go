// Package models — structs partagés entre db et handlers. Aucune logique ici, juste les types.
package models

import "time"

// User — un compte enregistré. PasswordHash stocké bcrypt, jamais le mdp en clair.
type User struct {
	ID           string
	Email        string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

// Session — cookie de session stocké en DB. Un seul par user (enforced dans db.CreateSession).
type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

// Category — sous-forum. Slugs pré-seedés au démarrage (voir db.seedCategories).
type Category struct {
	ID   int
	Name string
	Slug string
}

// Post — enrichi à la lecture : Username, Categories, Likes/Dislikes, UserReaction (-1/0/1).
type Post struct {
	ID           string
	UserID       string
	Username     string
	Title        string
	Content      string
	Categories   []Category
	Likes        int
	Dislikes     int
	UserReaction int
	CommentCount int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Comment — même logique que Post pour les compteurs de réactions.
type Comment struct {
	ID           string
	PostID       string
	UserID       string
	Username     string
	Content      string
	Likes        int
	Dislikes     int
	UserReaction int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TemplateData — données passées à chaque template. CurrentUser nil = visiteur anonyme.
type TemplateData struct {
	CurrentUser *User
	Posts       []Post
	Post        *Post
	Comments    []Comment
	Categories  []Category
	SelectedCat string
	Filter      string
	Error       string
	Success     string
	ErrCode     int
}
