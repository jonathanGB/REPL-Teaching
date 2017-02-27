# Liste des endpoints

* `GET /` : landing page, redirige vers `/users/signup` si pas identifié, sinon rediriges vers `/groups`


* `GET /users/signup` : page pour créer son compte
* `POST /users/signup` : soumettre formulaire de création de compte


* `GET /users/login` : page pour se connecter à son compte
* `POST /users/login` : soumettre les identifiants pour se connecter


* `GET /users/logout` : page de déconnexion
* `POST /users/logout` : se déconnecter


* `GET /groups` : voir tous les groupes liés à l'utilisateur
* `POST /groups` : créer un groupe, limité aux professeurs (via **AJAX**)


* `GET /groups/:groupId/join` : page pour joindre un groupe, redirige à `/groups/:groupId` si déjà dans le groupe ou vers `/groups` si c'est un professeur (peuvent pas joindre un groupe, seulement les étudiants)
* `POST /groups/:groupId/join` : joindre un groupe (via **AJAX**)


*Routes non-complétées*

* `GET /groups/:groupId` : dashboard pour un groupe précis, limité aux professeurs

* `GET /groups/:groupId/files` : voir tous les fichiers du groupe (visible pour l'utiliateur)


* `GET /groups/:groupId/files/:fileId` : vue d'éditeur du fichier (si l'étudiant n'a pas déjà ce fichier, redirigé vers le endpoint pour clôner)


* `POST /groups/:groupId/files/:fileId` : clôner le fichier **:fileId**


* `POST /groups/:groupId/files/` : créer un nouveau fichier


* `GET /groups/:groupId/files/:fileId/run` : rouler un fichier (ou méthode de la communication websockets?)
