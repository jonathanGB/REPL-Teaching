#!/bin/bash

# Before using this, two things are assumed!
# 1. You have Go installed, and have set your GOPATH
# 2. You have MongoDB installed, and have set the default path to data (/data/db)

sudo mongod > /dev/null &

for var in "$@"
do
	# to install Go dependencies to the project
	if [ "$var" = "--install" ]; then
		echo "Installing go-gin" && go get gopkg.in/gin-gonic/gin.v1 &&
		echo "Installing gin/multitemplate" && go get github.com/gin-contrib/multitemplate &&
		echo "Installing mgo.v2" && go get gopkg.in/mgo.v2 &&
		echo "Installing bson"   && go get gopkg.in/mgo.v2/bson &&
		echo "Installing bcrypt" && go get golang.org/x/crypto/bcrypt
	fi

	# to restart the database (creates by default a "foo" collection)
	if [ "$var" = "--restartDB" ]; then
		mongo repl --eval 'db.dropDatabase()' &&
		mongo repl --eval 'db.createCollection("foo"); db.foo.insert([{kek: true}, {kek: false}, {kek: true}])'
	fi
done


go run app.go
