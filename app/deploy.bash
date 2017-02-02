#!/bin/bash

# Before using this, two things are assumed!
# 1. You have Go installed, and have set your GOPATH
# 2. You have MongoDB installed, and have set the default path to data (/data/db)

for var in "$@"
do
	# to install Go dependencies to the project
	if [ "$var" = "--install" ]; then
		echo "Installing go-gin" && go get gopkg.in/gin-gonic/gin.v1 &&
		echo "Installing mgo.v2" && go get gopkg.in/mgo.v2 &&
		echo "Installing bson"   && go get gopkg.in/mgo.v2/bson
	fi

	# to restart the database (creates by default a "foo" collection)
	if [ "$var" = "--restartDB" ]; then
		mongo repl --eval 'db.dropDatabase()' &&
		mongo repl --eval 'db.createCollection("foo"); db.foo.insert([{kek: true}, {kek: false}, {kek: true}])'
	fi
done


sudo mongod > /dev/null &
go run app.go
