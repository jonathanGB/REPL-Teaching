# REPL Teaching - CSI 3540

### Un *Read-Eval-Print-Loop* (REPL) comme outil d'enseignment.

Deux types d'utilisateurs existent:

1. les professeurs
2. les étudiants.

#### Objectifs
Un professeur peut déployer pour sa classe des extraits (*snippets*) de code, que les étudiants pourront ensuite cloner pour modifier individuellement. Le professeur peut voir un tableau de bord qui montrera les changements et activités des étudiants en temps réel, que le professeur peut utiliser comme mesure de participation des étudiants ou pour des fins de discussion devant la classe. Évidemment, les extraits de code seront exécutables directement dans le navigateur (REPL), avec un éditeur de texte adapté pour du code.

#### Technologies
* HTML & CSS
* JavaScript && jQuery
* Go (serveur) + Micro-framework Gin
* WebSockets (temps-réel)
* MongoDB (base de données)
* Compilateurs / interpréteurs des langages qui seront couverts par le REPL


#### Comment déployer
Pour déployer localement, il est nécessaire d'avoir **Go** (`GOPATH` bien configuré aussi) et **MongoDB** d'installés (le dossier par défaut pour les données doit être utilisé `/data/db`). Un script `deploy.bash` permet ensuite d'installer automatiquement les dépendances **go** du projet, configurer la base de données et lancer le tout. Ce script est écrit pour les sytèmes **UN*X**.

Certains *flags* sont disponibles pour lancer le script, les voici:
* `--install`: installer les dépendances **go**
* `--restartDB`: réinitialiser la base de données

#### Comment tester
À partir de la ligne de commande, dirigez-vous dans le dossier **app**, puis simplement exécuter `go test`, et les tests rouleront automatiquement.
