package db

import "database/sql"

// ResultatReaction contient les compteurs après un vote
type ResultatReaction struct {
	Likes        int
	Dislikes     int
	UserReaction int // 1 = like, -1 = dislike, 0 = aucun vote
}

// ToggleReaction ajoute, change ou annule un vote (like/dislike)
func (d *DB) ToggleReaction(idUtilisateur, idCible, typeCible string, valeur int) (*ResultatReaction, error) {
	var voteExistant int
	err := d.conn.QueryRow(
		`SELECT value FROM reactions WHERE user_id = ? AND target_id = ? AND target_type = ?`,
		idUtilisateur, idCible, typeCible,
	).Scan(&voteExistant)

	if err == sql.ErrNoRows {
		// pas encore voté → on ajoute
		_, err = d.conn.Exec(
			`INSERT INTO reactions (user_id, target_id, target_type, value) VALUES (?, ?, ?, ?)`,
			idUtilisateur, idCible, typeCible, valeur,
		)
	} else if err == nil {
		if voteExistant == valeur {
			// même vote → on annule
			_, err = d.conn.Exec(
				`DELETE FROM reactions WHERE user_id = ? AND target_id = ? AND target_type = ?`,
				idUtilisateur, idCible, typeCible,
			)
			valeur = 0
		} else {
			// vote différent → on change
			_, err = d.conn.Exec(
				`UPDATE reactions SET value = ? WHERE user_id = ? AND target_id = ? AND target_type = ?`,
				valeur, idUtilisateur, idCible, typeCible,
			)
		}
	}

	if err != nil {
		return nil, err
	}

	// on retourne les nouveaux compteurs
	resultat := &ResultatReaction{UserReaction: valeur}
	d.conn.QueryRow(`
		SELECT
			COALESCE(SUM(CASE WHEN value =  1 THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN value = -1 THEN 1 ELSE 0 END), 0)
		FROM reactions WHERE target_id = ? AND target_type = ?`,
		idCible, typeCible,
	).Scan(&resultat.Likes, &resultat.Dislikes)

	return resultat, nil
}
