package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"

	forum "forumJS/forumJS"
)

func main() {
	// forum.HashPassword("abc")
	// forum.HashPassword("abc")
	// Les méthode HTTP :
	http.HandleFunc("/Inscription", func(w http.ResponseWriter, r *http.Request) {
		Inscription(w, r)
	})

	http.HandleFunc("/Connexion", func(w http.ResponseWriter, r *http.Request) {
		Connexion(w, r)
	})

	http.HandleFunc("/ChangerPage", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		pageDemandé := r.FormValue("pageDemandé")
		http.ServeFile(w, r, "pages/"+pageDemandé)
	})

	http.HandleFunc("/ChangerPageSpéciale", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		pageDemandé := r.FormValue("pageDemandé")
		forum.RevenirSurLaPageAccueil(w, r, 0, false, true, 0, pageDemandé);
	})

	http.HandleFunc("/InteractionPost", func(w http.ResponseWriter, r *http.Request) {
		forum.InteractionPost(w, r)
	})

	http.HandleFunc("/EnvoyerCommentaire", func(w http.ResponseWriter, r *http.Request) {
		EnvoyerCommentaire(w, r)
	})

	http.HandleFunc("/AjouterEspaceCommentaire", func(w http.ResponseWriter, r *http.Request) {
		forum.AjouterEspaceCommentaire(w, r)
	})

	http.HandleFunc("/ChangerDeFilDeDiscution", func(w http.ResponseWriter, r *http.Request) {
		forum.ChangerDeFilDeDiscution(w, r)
	})

	http.HandleFunc("/PartagerPage", func(w http.ResponseWriter, r *http.Request) {
		PartagerPage(w, r)
	})

	http.HandleFunc("/NouveauFilDeDiscution", func(w http.ResponseWriter, r *http.Request) {
		forum.NouveauFilDeDiscution(w, r)
	})

	http.HandleFunc("/BarreDeRecherche", func(w http.ResponseWriter, r *http.Request) {
		forum.Recherche(r.FormValue("Recherche"), w, r)
	})

	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("./style"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	http.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("./pages"))))

	fmt.Println(forum.CheckPassword(forum.HashPassword("abc"), "abc"))
	// Au démarage du serveur :
	log.Println("Serveur lancé sur http://localhost:8080")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		iD_publication_commentaire := r.FormValue("iD_publication_commentaire")
		valeur_iD_publication_commentaire := -1
		if iD_publication_commentaire != "" {
			valeur, err := strconv.Atoi(iD_publication_commentaire)
			if err == nil {
				valeur_iD_publication_commentaire = valeur
			}
		}

		valeur := (r.FormValue("iD_fil_de_discussion"))
		iD_fil_de_discussion, err := strconv.Atoi(valeur)
		if err != nil {
			if r.FormValue("PageSpéciale") == "TableauDeBord" {
				forum.TableauDeBord(w, r)
			}
			forum.ComplétéLaPageAccueil(w, r)
		} else {
			// fmt.Println("on est en : ",iD_fil_de_discussion)
			forum.ComplétéLaPageForum(w, r)
			forum.AfficherToutLesPost(iD_fil_de_discussion, w, r, valeur_iD_publication_commentaire)
		}
	})

	forum.InitDB()

	http.HandleFunc("/open", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		go func() {
			_ = exec.Command("xdg-open", "http://localhost:8080/").Start()
		}()
		w.Write([]byte("Attempted to open browser"))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func Inscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// w.Write([]byte(DémarerUnePartie(informations, r)))
	email := r.FormValue("email")
	password := r.FormValue("password")
	nomUtilisateur := r.FormValue("nomUtilisateur")

	réusie := forum.AjouterUnUtilisateur(w, email, password, nomUtilisateur)
	fmt.Println(réusie)
	if réusie {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.ServeFile(w, r, "pages/inscription.html")
	}
}

func Connexion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// w.Write([]byte(DémarerUnePartie(informations, r)))
	email := r.FormValue("email")
	password := r.FormValue("password")

	réusie := false
	iD_Utilisateur := forum.ConnecterUtilisateur(email, password)
	if iD_Utilisateur != 0 {
		forum.CrééUnCookie(w, iD_Utilisateur)
		réusie = true
	}
	fmt.Println(réusie)

	if réusie {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.ServeFile(w, r, "pages/inscription.html")
	}
}

func EnvoyerCommentaire(w http.ResponseWriter, r *http.Request) {
	valeur := (r.FormValue("answer"))
	answer, err := strconv.Atoi(valeur)
	if err != nil {
		fmt.Println(err)
		answer = 0
	}

	valeur = (r.FormValue("iD_fil_de_discussion"))
	fmt.Println("form : ", valeur)
	iD_fil_de_discussion, err := strconv.Atoi(valeur)
	if err != nil {
		iD_fil_de_discussion = 0
	}
	fmt.Println(iD_fil_de_discussion)

	annulerCommentaire := r.FormValue("AnnulerCommentaire")
	if annulerCommentaire == "oui" {
		forum.RevenirSurLaPageAccueil(w, r, answer, true, true, -1, "")
		return
	}

	leTexte := r.FormValue("leTexte")
	idUtilisateur := forum.VérifierCookie(r)
	if idUtilisateur > 0 {
		threadID := iD_fil_de_discussion
		dsnURI := "db/forum.db"
		db, err := sql.Open("sqlite", dsnURI)
		if err != nil {
			fmt.Println("Erreur d'ouverture :", err)
			return
		}

		defer db.Close()

		forum.CreatePost(idUtilisateur, threadID, leTexte, db, answer)
	}

	forum.RevenirSurLaPageAccueil(w, r, answer, false, false, -1, "")

}

func PartagerPage(w http.ResponseWriter, r *http.Request) string {
	valeur := (r.FormValue("iD_fil_de_discussion"))
	iD_fil_de_discussion, err := strconv.Atoi(valeur)
	if err != nil {
		fmt.Println(err)
		iD_fil_de_discussion = 0
	}

	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}

	if pos := strings.Index(referer, "?"); pos != -1 {
		referer = referer[:pos]
	}

	referer = fmt.Sprintf("%s", referer)
	referer += fmt.Sprintf("iD_fil_de_discussion=%d", iD_fil_de_discussion)
	return referer
}
