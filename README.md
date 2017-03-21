# REPL Teaching - CSI 3540

### Un *Read-Eval-Print-Loop* (REPL) comme outil d'enseignment.

Deux types d'utilisateurs existent:

1. les professeurs
2. les étudiants.

#### Objectifs
Un professeur peut déployer pour sa classe des extraits (*snippets*) de code, que les étudiants pourront ensuite cloner pour modifier individuellement. Le professeur peut voir un tableau de bord qui montrera les changements et activités des étudiants en temps réel, que le professeur peut utiliser comme mesure de participation des étudiants ou pour des fins de discussion devant la classe. Évidemment, les extraits de code seront exécutables directement dans le navigateur (REPL), avec un éditeur de texte adapté pour du code.

#### Technologies
* HTML5 & CSS3
	* testé sur les versions récentes stables de Chrome + Firefox
* JavaScript & jQuery
	* j'utilise des fonctionnalités de ES6
* Go (serveur) + Micro-framework Gin
	* nécessite Go `1.7+` (je roule `1.8`, la version la plus récente)
* WebSockets (temps-réel)
* MongoDB (base de données)
	* je roule `3.2.7`, je recommande d'avoir au moins `3.2.x`
* Docker (containers pour exécuter les scripts du REPL)
	* dernière version stable sur macOS (`17.03`)
* BASH (exécuter les scripts de déploiement)
	* nécessite au moins `4.3` pour rouler le script de détection d'images docker automatique. Il peut être utile d'utiliser le flag `--noBuild` dans le **deploy** pour ignorer ce script.

#### Langages supportés par le REPL
* JavaScript (via Node)
* Go
* Python
* Ruby
* Elixir

#### Comment déployer
Pour déployer localement, il est nécessaire d'avoir **Go** (`GOPATH` bien configuré aussi) et **MongoDB** d'installés (le dossier par défaut pour les données doit être utilisé `/data/db`). Un script `deploy.bash` permet ensuite d'installer automatiquement les dépendances **go** du projet, configurer la base de données, bâtir les images docker manquantes, et lancer le tout. Ce script est écrit pour les sytèmes **UN*X**.

##### À mentionner que ce script doit être lancé à partir du même dossier, sinon les références relatives vont être disfonctionnelles. (i.e. lancer `./deploy.bash` et non `./app/deploy.bash`)

Si ce n'est pas déjà fait, il faut s'assurer que l'utilisateur lançant le script de déploiement est le "owner" du dossier des données `/data/db`. Si ce n'est pas le cas, simplement exécuter `chown -R <USER> /data`.

Certains *flags* sont disponibles pour lancer le script, les voici:
* `--install`: installer les dépendances **go**
* `--restartDB`: réinitialiser la base de données
* `--noBuild`: ne pas lancer le script de détection et d'installation automatique des images docker manquantes (ce script nécessite BASH `4.3`)


#### Comment tester
À partir de la ligne de commande, dirigez-vous dans le dossier **app**, puis simplement exécuter `go test`, et les tests rouleront automatiquement.
