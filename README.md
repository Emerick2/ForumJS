
```
db
	forum.db	                            # Basses de données contenant les messages.		
	interactionUtilisateur.db               # Base de données contenants les intéractions des utilisateurs.
	threads.db                              # Basse de données contenant la liste des sujet de discution.
	user.db                                 # Basse de données contenant les utilisateur

forumJS
	AfficherLesPostes.go                    # Affichage les postes
	compétéPageAccueil.go                   # Complété la page d'accueil.
	cookie.go                               # Créé et vérifier les cookies de secions.
	CrééUnFilDeDiscution.go                 # Créé un fil de discution.
	database_Post.go                        # Sauvegarder les postes dans la base de données.
	date.go                                 # Transformer une date Time.Time en un string bien formater.
	deconnexion.go                          # Permettre à l'utilisateur de se déconnecter.
	FilDeDiscution.go                       # Gestion des fils de discutions.
	gestionUtilisateur.go                   # Gestion des utilisateurs.
	http_errors.go                          # Gestion des erreur http.
	motdepasse.go                           # Gestion des mots de passes.
	post.go                                 # Sauvegarde, lecture et structure des postes.
	postes.go                               # Gestion de l'intéraction des utilisateurs avec les postes.
	Recherche.go                            # Barre de recherche.
	RevenirSurLaPageAccueilScript.go        # Script pour changer de pages dynamiquement.
	StructureUtilisateur.go                 # Structure des utilisateurs.
	TableauDeBord.go                        # Complété la page tableau de bord.

images                                      # Un dossier contenant toutes les images du projet.

js
	Partager.js                             # Enregistrer le lien de la page dans le press papier de l'utilisateur.

pages
	discution.html                          # Page ou l'on peut lire les postes.
	inscription.html                        # Page pour s'inscrire.
	main.html                               # Page d'accueil.
	nouveau-sujet.html                      # Page pour écrire un nouveau sujet.
	tableau-de-bord.html                    # Page du tableau de bord.
	template-commentaire.html               # Template pour ajouter un commenataire.
	template-fiche-utilisateur.html         # Template pour ajouter une fiche d'utilisateur.
	template-haut-file.html                 # Template pour ajouter les premières partie de fil de discution.
	template-post.html                      # Template pour les posts.

style
	barre-de-recherche.css                  # Barre de recherche.
	bloc-profile.css                        # Partie des profiles.
	commentaire.css                         # Partie commentaire.
	NouveauFilDeDiscution.css               # Fil de discution.
	pageAccueil.css                         # Page d'accueil.
	pageConnexion.css                       # Page de connexion.
	post.css                                # Partie des postes.
	profile.css                             # Partie des profiles.
	tableau-de-bord.css                     # Page du tableau de bord.

main.go                                     # C'est le script qui lance le programme.
```
