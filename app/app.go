package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"fmt"
)

func getMainEngine() *gin.Engine {
	router := gin.Default()
	router.Static("/public", "./public")
	router.StaticFile("/favicon.ico", "./public/img/favicon.ico")

	s := initDB()
	defer s.Close()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.GET("/foo", func(c *gin.Context) {
		fooObjs := []struct {
			Id  bson.ObjectId `json:"id" bson:"_id"`
			Kek bool          `json:"topkek" bson:"kek"`
		}{}

		kekFilter, err := strconv.ParseBool(c.Query("kek"))

		var queryFilter bson.M
		if err == nil {
			queryFilter = bson.M{
				"kek": kekFilter,
			}
		}

		s.DB("repl").C("foo").Find(queryFilter).All(&fooObjs)
		c.JSON(http.StatusOK, fooObjs)
	})

	fmt.Println("\n") // empty buffer in output
	return router
}

func initDB() *mgo.Session {
	s, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}

	return s
}

func main() {
	getMainEngine().Run(":8080")
}
