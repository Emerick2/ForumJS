package handlers

import (
	"encoding/json"
	"net/http"
)

// requeteReaction représente ce que le JavaScript envoie quand on clique sur like/dislike
type requeteReaction struct {
	IdCible    string `json:"target_id"`
	TypeCible  string `json:"target_type"` // "post" ou "comment"
	Valeur     int    `json:"value"`       // 1 pour like, -1 pour dislike
}

// handleReaction traite un like ou dislike envoyé par le JavaScript
func (app *App) handleReaction(w http.ResponseWriter, r *http.Request) {
	utilisateur := app.utilisateurConnecte(r)

	// on lit le JSON envoyé par le navigateur
	var req requeteReaction
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, `{"error":"données invalides"}`, http.StatusBadRequest)
		return
	}

	// on vérifie que le type est valide
	if req.TypeCible != "post" && req.TypeCible != "comment" {
		http.Error(w, `{"error":"type invalide"}`, http.StatusBadRequest)
		return
	}

	// on vérifie que la valeur est 1 ou -1
	if req.Valeur != 1 && req.Valeur != -1 {
		http.Error(w, `{"error":"valeur invalide"}`, http.StatusBadRequest)
		return
	}

	// on enregistre la réaction dans la base de données
	// si l'utilisateur a déjà voté pareil, ça annule son vote
	// si l'utilisateur a voté différemment, ça change son vote
	resultat, err := app.db.ToggleReaction(utilisateur.ID, req.IdCible, req.TypeCible, req.Valeur)
	if err != nil {
		http.Error(w, `{"error":"erreur serveur"}`, http.StatusInternalServerError)
		return
	}

	// on répond avec les nouveaux compteurs en JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"likes":         resultat.Likes,
		"dislikes":      resultat.Dislikes,
		"user_reaction": resultat.UserReaction,
	})
}
