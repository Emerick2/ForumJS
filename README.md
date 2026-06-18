# FORUSKY - le forum

### Contexte du Projet

L'objectif de notre projet était de réaliser un forum en utilisant les langages de programmation Go, SQL, HTML, CSS et JS.
Nous avons choisi le nom Forusky puisqu'il s'agit d'une fusion entre Husky (une race de chien) et forum.

### Structure du Projet

```
db
    forum.db                                # Base de données contenant les messages.     
    interactionUtilisateur.db               # Base de données contenant les interactions des utilisateurs.
    threads.db                              # Base de données contenant la liste des sujets de discussion.
    user.db                                 # Base de données contenant les utilisateurs.

forumJS
    AfficherLesPostes.go                    # Affichage des posts.
    compétéPageAccueil.go                   # Complète la page d'accueil.
    cookie.go                               # Crée et vérifie les cookies de session.
    CrééUnFilDeDiscution.go                 # Crée un fil de discussion.
    database_Post.go                        # Sauvegarde les posts dans la base de données.
    date.go                                 # Transforme une date time.Time en une chaîne bien formatée.
    deconnexion.go                          # Permet à l'utilisateur de se déconnecter.
    FilDeDiscution.go                       # Gestion des fils de discussion.
    gestionUtilisateur.go                   # Gestion des utilisateurs.
    http_errors.go                          # Gestion des erreurs HTTP.
    motdepasse.go                           # Gestion des mots de passe.
    post.go                                 # Sauvegarde, lecture et structure des posts.
    postes.go                               # Gestion de l'interaction des utilisateurs avec les posts.
    Recherche.go                            # Barre de recherche.
    RevenirSurLaPageAccueilScript.go        # Script pour changer de page dynamiquement.
    StructureUtilisateur.go                 # Structure des utilisateurs.
    TableauDeBord.go                        # Complète la page du tableau de bord.

images                                      # Dossier contenant toutes les images du projet.

js
    Partager.js                             # Enregistre le lien de la page dans le presse-papiers de l'utilisateur.

pages
    discution.html                          # Page où l'on peut lire les posts.
    inscription.html                        # Page pour s'inscrire.
    main.html                               # Page d'accueil.
    nouveau-sujet.html                      # Page pour écrire un nouveau sujet.
    tableau-de-bord.html                    # Page du tableau de bord.
    template-commentaire.html               # Template pour ajouter un commentaire.
    template-fiche-utilisateur.html         # Template pour ajouter une fiche d'utilisateur.
    template-haut-file.html                 # Template pour ajouter la première partie d'un fil de discussion.
    template-post.html                      # Template pour les posts.

style
    barre-de-recherche.css                  # Style de la barre de recherche.
    bloc-profile.css                        # Style de la partie des profils.
    commentaire.css                         # Style de la partie commentaires.
    NouveauFilDeDiscution.css               # Style du fil de discussion.
    pageAccueil.css                         # Style de la page d'accueil.
    pageConnexion.css                       # Style de la page de connexion.
    post.css                                # Style de la partie des posts.
    profile.css                             # Style de la partie des profils.
    tableau-de-bord.css                     # Style de la page du tableau de bord.

main.go                                     # Script principal qui lance le programme.

```

### Technologies Utilisées

* **Backend :** Go, SQL
* **Frontend :** HTML5, CSS3, JavaScript

### Comment l'utiliser ?

```bash
go run .

```

### Réalisé par :

* Émerick
* Benjamin
* Paul Elie