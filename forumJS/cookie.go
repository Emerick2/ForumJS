package forumjs

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var cleSecrete = []byte("shjzaqfjkffzf5ver6ezrcf8rceez569z")

func CrééUnCookie(w http.ResponseWriter, id int) {
	durée := 3 * time.Hour
	expiration := time.Now().Add(durée)

	idTexte := strconv.Itoa(id)

	mac := hmac.New(sha256.New, cleSecrete)
	mac.Write([]byte(idTexte))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	valeurSecurisée := idTexte + "." + signature

	cookie := &http.Cookie{
		Name:     "session_utilisateur",
		Value:    valeurSecurisée,
		Expires:  expiration,
		MaxAge:   int(durée.Seconds()),
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	}

	http.SetCookie(w, cookie)
}

func VérifierCookie(r *http.Request) int {
	cookie, err := r.Cookie("session_utilisateur")
	if err != nil {
		fmt.Println("Attention, l'utilisateur n'est pas connecté !")
		return 0
	}

	valeur := cookie.Value

	var idTexte, signatureRecue string
	for i := 0; i < len(valeur); i++ {
		if valeur[i] == '.' {
			idTexte = valeur[:i]
			signatureRecue = valeur[i+1:]
			break
		}
	}

	if idTexte == "" || signatureRecue == "" {
		fmt.Println("Attention, l'utilisateur n'est pas connecté !")
		return 0
	}

	mac := hmac.New(sha256.New, cleSecrete)
	mac.Write([]byte(idTexte))
	signatureAttendue := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if signatureRecue != signatureAttendue {
		fmt.Println("Tentative de modification de cookie détectée !")
		return 0
	}

	id, err := strconv.Atoi(idTexte)
	if err != nil {
		fmt.Println("Attention, l'utilisateur n'est pas connecté !")
		return 0
	}

	return id
}
