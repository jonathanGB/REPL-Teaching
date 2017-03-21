#!/usr/bin/env bash

# Before using this, two things are assumed!
# 1. You have Go installed, and have set your GOPATH
# 2. You have MongoDB installed, and have set the default path to data (/data/db)

buildDockerImagesFlag="true"

mongod > /dev/null &

for var in "$@"
do
	# to install Go dependencies to the project
	if [ "$var" = "--install" ]; then
		echo "Installing go-gin" && go get gopkg.in/gin-gonic/gin.v1 &&
		echo "Installing gin/multitemplate" && go get github.com/gin-contrib/multitemplate &&
		echo "Installing mgo.v2" && go get gopkg.in/mgo.v2 &&
		echo "Installing bson"   && go get gopkg.in/mgo.v2/bson &&
		echo "Installing bcrypt" && go get golang.org/x/crypto/bcrypt &&
		echo "Installing jwt" && go get github.com/dgrijalva/jwt-go &&
		echo "Installing gorilla-toolkit/websocket" && go get github.com/gorilla/websocket &&
		echo "Installing docker/api/types" && go get "github.com/docker/docker/api/types" &&
		echo "Installing docker/api/types/container" && go get "github.com/docker/docker/api/types/container" &&
		echo "Installing docker/client" && go get "github.com/docker/docker/client" &&
		echo ""
	fi

	# to restart the database (creates by default a "foo" collection)
	if [ "$var" = "--restartDB" ]; then
		mongo repl --eval 'db.dropDatabase()' &&
		mongo repl --eval 'db.createCollection("foo"); db.foo.insert([{kek: true}, {kek: false}, {kek: true}])' #&&
		#mongo repl --eval 'db.createCollection("users"); db.users.insert([{ "_id" : ObjectId("58ae488a9d1fa299c7dc795a"), "name" : "topprof", "email" : "foo@bar", "role" : "teacher", "groups" : [ObjectId("58ae4b00462a70ed0c38367d") ], "password" : BinData(0,"JDJhJDEwJENwYmx3elNiWGhiVlpFWVZyUFdlY090VHJuQW0vbjdVc1pDMFpuQUpJQWhoaWQ0RHV1WXUy") }, { "_id" : ObjectId("58ae48e39d1fa299c7dc795b"), "name" : "topstudent", "email" : "student@bar", "role" : "student", "groups" : [ ObjectId("58ae4b00462a70ed0c38367d") ], "password" : BinData(0,"JDJhJDEwJGVqMk5Nd0JXbUJhRDlCalpBdXR5dnVnZmU5NlN0LkZQY0J3MS5VUlNRVFppWW9DLjhGOHhL") }]);' &&
		#mongo repl --eval 'db.createCollection("groups"); db.groups.insert({"_id": ObjectId("58ae4b00462a70ed0c38367d"), name: "fou groupe", teacher: ObjectId("58ae488a9d1fa299c7dc795a"), "teacher-name": "topprof", files: [], password: BinData(0,"JDJhJDEwJENwYmx3elNiWGhiVlpFWVZyUFdlY090VHJuQW0vbjdVc1pDMFpuQUpJQWhoaWQ0RHV1WXUy")});'
	fi

	#cancel docker Images build if not found default behaviour
	if [ "$var" = "--noBuild" ]; then
		buildDockerImagesFlag="false"
	fi
done

if [ "$buildDockerImagesFlag" == "true" ]; then
	./buildDockerImages.bash
fi


go run app.go
