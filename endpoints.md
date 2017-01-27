# Liste des endpoints

* `GET /` : landing page, liens vers **signup** et **login**
* `GET /signup` : page pour créer son compte
* `GET /login` : page pour se connecter à son compte
* `POST /logout` : se déconnecter
* `GET /groups` : voir tous les groupes liés à l'utilisateur
* `GET /groups/:groupId` : dashboard pour un groupe précis
* `GET /groups/:groupId/files` : voir tous les fichiers du groupe (visible pour l'utiliateur)
* `GET /groups/:groupId/files/:fileId` : vue d'éditeur du fichier (si l'étudiant n'a pas déjà ce fichier, redirigé vers le endpoint pour clôner)
* `POST /groups/:groupId/files/:fileId` : clôner le fichier :fileId
* `POST /groups/:groupId/files/` : créer un nouveau fichier
* `GET /groups/:groupId/files/:fileId/run` : rouler un fichier (ou méthode de la communication websockets?)
