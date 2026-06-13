package forumjs

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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

	// referer = ""
	// listePartieURL := strings.Split(referer, "/")
	// fmt.Println(listePartieURL)
	// for i := 0; i < len(listePartieURL)-1; i++ {
	// 	referer += listePartieURL[i];
	// 	if (i+1 < len(listePartieURL)-1){
	// 		referer += "/"
	// 	}
	// }

	// fmt.Println("url : ",referer)
	// if pos := strings.Index(referer, "?"); pos != -1 {
	// 	referer = referer[:pos]
	// }

	// if strings.Contains(referer, "/BarreDeRecherche") {
	// 	referer = strings.Replace(referer, "/BarreDeRecherche", "/", 1)
	// }

	u, err := url.Parse(referer)
	if err != nil {
		return
	}

	referer = u.Scheme + "://" + u.Host + "/"

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
