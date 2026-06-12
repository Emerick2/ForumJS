package forumjs

import (
	"fmt"
	"net/http"
)

func Recherche(recherche string, w http.ResponseWriter, r *http.Request) {
	listeThread, err := GetThread()
	if err != nil {
		fmt.Println(err)
	}

	listePosts := GetPost()

	textePris := []string{}
	textePrisFil := []string{}
	nombreRésultaFil := []Thread{}
	nombreRésultaMessage := []Post{}

	for i := 0; i < len(listeThread); i++ {
		mot := listeThread[i].Message_content
		peutÊtreVuAvecSeTermeDeRecherche := !EstDansLaListe(mot, textePris) && PeutÊtreVuAvecSeTermeDeRecherche(mot, recherche)
		if peutÊtreVuAvecSeTermeDeRecherche {
			nombreRésultaFil = append(nombreRésultaFil, listeThread[i])
			textePris = append(textePris, mot)
			textePrisFil = append(textePris, mot)
		}
	}

	for i := 0; i < len(listePosts); i++ {
		mot := listePosts[i].Content
		peutÊtreVuAvecSeTermeDeRecherche := !EstDansLaListe(mot, textePris) && PeutÊtreVuAvecSeTermeDeRecherche(mot, recherche)
		if peutÊtreVuAvecSeTermeDeRecherche {
			nombreRésultaMessage = append(nombreRésultaMessage, listePosts[i])
			textePris = append(textePris, mot)
		} else if EstDansLaListe(mot, textePrisFil) {
			nouvelleListe := []Post{listePosts[i]}
			nombreRésultaMessage = append(nouvelleListe, nombreRésultaMessage...)
		}
	}
	AfficherRecherche(w, r, nombreRésultaMessage)
}

func EstDansLaListe(mot string, tableau []string) bool {
	for i := 0; i < len(tableau); i++ {
		if mot == tableau[i] {
			return true
		}
	}
	return false
}

/*
La recherche sur se forum permetteras à l'utilisateur d'écrire se qu'il cherche
exemple "livre de romain"

On vas trié par pertinance.
Le plus pertinant dans l'ordre est :
- le mot clef est dans le titre ?
- Le fil est aimé et commenter ?
- le fil est récent ?

- Le mot clef est dans les commentaire ?

Pour avoir le plus de pertinance possible, je vais :
1 - télécharger tous les fils de discution
2 - utiliser l'algorythme de recherche par mot clef pour ne garder que ceux qui ont le mot clef chercher.
3 - S'il y en à 5 ou plus :
	- les trié par pertinances
	- FINI !

4 - SINON :
	5 - télécharger tous les messages
	6 - utiliser l'algorythme de recherche par mot clef pour ne garder que ceux qui ont le mot clef chercher.
	7 - Supprimer tous les doubles (qui sont apparus à la fois dans le titre et la description)
	8 - les trié par pertinances
	- FINI !
*/

/*
Que doit faire la barre de recherche ?
1 - donner se qui se rapprocher à 3 lettre prait du mots chercher
cherche : poisson
trouve : poison, poisson, poisons, poissonier
2 - Donne se qui contient la chose chercher à l'intérieur :
cherche : qui
trouve : quimange, avecquiilest, qui
*/

func PeutÊtreVuAvecSeTermeDeRecherche(résultat string, recherche string) bool {
	if résultat == "" || recherche == "" {
		return false
	}
	résultat = ToUpper(résultat)
	recherche = ToUpper(recherche)
	// 	recherche : 'qui'
	// 	résultat : 'qui mange', 'avec qui il est', 'qui'
	if résultat == recherche || len(recherche) == 0 {
		return true
	}
	max := len(résultat)
	compte := 0
	if len(recherche) <= max {
		for i := 0; i < len(résultat); i++ {
			if compte < len(recherche) {
				if résultat[i] == recherche[compte] {
					compte++
					if compte >= len(recherche) {
						return true
					}
				} else {
					compte = 0
				}
			}
		}
	}
	return false
}

func ToUpper(texte string) string {
	//cette fonction ne fonctionne pas sur tout les accents.
	résultat := ""
	if texte == "" {
		return résultat
	}
	runes := []rune(texte)
	for i := 0; i < len(texte); i++ {
		if runes[i] >= 97 && runes[i] <= 122 {
			résultat += (string)(runes[i] - 32)
		} else {
			listeMinuscule := []rune{'é', 'è', 'ô', 'û', 'â', 'ê', 'î', 'ö', 'ë', 'ü', 32}
			listeMajuscule := []rune{'É', 'È', 'Ô', 'Û', 'Â', 'Ê', 'Î', 'Ö', 'Ë', 'Ü', 0}
			vu := false
			for j := 0; j < len(listeMinuscule); j++ {
				if len(listeMajuscule) > j {
					if runes[i] == listeMinuscule[j] {
						runes[i] = listeMajuscule[j]
						vu = true
						break
					}
				}
			}
			if !vu {
				résultat += (string)(runes[i])
			}
		}
	}
	fmt.Println(résultat)
	return résultat
}

func AfficherRecherche(w http.ResponseWriter, r *http.Request, nombreRésultaMessage []Post) {
	ComplétéLaPageForum(w, r)

	for i := 0; i < len(nombreRésultaMessage); i++ {
		AfficherPost(nombreRésultaMessage[i], w, r, false, 0, false)
	}
}
