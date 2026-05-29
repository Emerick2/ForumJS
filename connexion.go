package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"

	_ "modernc.org/sqlite"

	forum "forumJS/forumJS"
)

func main() {
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

	http.HandleFunc("/InteractionPost", func(w http.ResponseWriter, r *http.Request) {
		InteractionPost(w, r)
	})

	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("./style"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	http.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("./pages"))))

	// Au démarage du serveur :
	log.Println("Serveur lancé sur http://localhost:8080")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		forum.ComplétéLaPageAccueil(w, r)
		forum.AfficherPost(0, w, r)
	})

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

func InteractionPost(w http.ResponseWriter, r *http.Request) {
	nomAction := r.FormValue("nomAction")
	iD_publication := r.FormValue("iD_publication")
	iD_fil_de_discussion := r.FormValue("iD_fil_de_discussion")

	fmt.Println(nomAction)
	fmt.Println(iD_publication)
	fmt.Println(iD_fil_de_discussion)
}
