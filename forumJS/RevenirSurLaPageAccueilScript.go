package forumjs

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func RevenirSurLaPageAccueil(w http.ResponseWriter, r *http.Request, iD_publication int, changerCommentaire bool) {
	// referer := r.Header.Get("Referer")
	// if referer == "" {
	// 	referer = "/"
	// }
	// if iD_publication > 0 {
	// 	referer = fmt.Sprintf("%s#post-%d", referer, iD_publication)
	// }
	// http.Redirect(w, r, referer, http.StatusSeeOther)

	valeur := (r.FormValue("iD_fil_de_discussion"))
	iD_fil_de_discussion, err := strconv.Atoi(valeur)
	if err != nil {
		iD_fil_de_discussion = 0
	}

	valeur = (r.FormValue("iD_publication_commentaire"))
	iD_publication_commentaire, err := strconv.Atoi(valeur)
	if err != nil {
		iD_publication_commentaire = -1
	}
	if changerCommentaire {
		iD_publication_commentaire = iD_publication
	}

	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}

	if pos := strings.Index(referer, "?"); pos != -1 {
		referer = referer[:pos]
	}

	if iD_publication > 0 {
		referer = fmt.Sprintf("%s?iD_publication_commentaire=%d&iD_fil_de_discussion=%d#post-%d", referer, iD_publication_commentaire, iD_fil_de_discussion, iD_publication)
	}

	http.Redirect(w, r, referer, http.StatusSeeOther)
}
