package forumjs

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func RevenirSurLaPageAccueil(w http.ResponseWriter, r *http.Request, iD_publication int, changerCommentaire bool, nePlusSélectionnerUnCommentaire bool, nouveauFilDeDiscution int, pageSpéciale string) {
	valeur := (r.FormValue("iD_fil_de_discussion"))
	iD_fil_de_discussion, err := strconv.Atoi(valeur)
	if err != nil {
		iD_fil_de_discussion = 0
	}
	if nouveauFilDeDiscution != -1 {
		iD_fil_de_discussion = nouveauFilDeDiscution
	}

	valeur = (r.FormValue("iD_publication_commentaire"))
	iD_publication_commentaire, err := strconv.Atoi(valeur)
	if err != nil {
		iD_publication_commentaire = -1
	}
	if changerCommentaire {
		iD_publication_commentaire = iD_publication
	}
	if nePlusSélectionnerUnCommentaire {
		iD_publication_commentaire = -1
	}

	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
		// referer = "/discution.html"
	}

	if pos := strings.Index(referer, "?"); pos != -1 {
		referer = referer[:pos]
	}

	nombreAjout := 0
	referer = fmt.Sprintf("%s", referer)
	if pageSpéciale == "" {
		if iD_publication > 0 {
			referer += TernaireStr(nombreAjout > 0, "&", "?")
			nombreAjout++
			referer += fmt.Sprintf("iD_publication_commentaire=%d", iD_publication_commentaire)
		}
		if iD_fil_de_discussion >= 0 {
			referer += TernaireStr(nombreAjout > 0, "&", "?")
			nombreAjout++
			referer += fmt.Sprintf("iD_fil_de_discussion=%d", iD_fil_de_discussion)
		}
	} else {
		referer += TernaireStr(nombreAjout > 0, "&", "?")
		nombreAjout++
		referer += fmt.Sprintf("PageSpéciale=%s", pageSpéciale)
	}
	if iD_publication > 0 {
		nombreAjout++
		referer += fmt.Sprintf("#post-%d", iD_publication)
	}

	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func TernaireStr(condition bool, valeur1 string, valeur2 string) string {
	if condition {
		return valeur1
	} else {
		return valeur2
	}
}