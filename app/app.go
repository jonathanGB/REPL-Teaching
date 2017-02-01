package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func main() {
	router := gin.Default()
	router.Static("/public", "./public")
	router.StaticFile("/favicon.ico", "./public/img/favicon.ico")

	s := initDB()
	defer s.Close()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.GET("/foo", func(c *gin.Context) {
		fooObj := struct {
			Id  bson.ObjectId `json:"id" bson:"_id"`
			Kek bool          `json:"topkek" bson:"kek"`
		}{}

		s.DB("repl").C("foo").Find(bson.M{}).One(&fooObj)
		c.JSON(http.StatusOK, fooObj)
	})

	router.Run(":8080")
}

func initDB() *mgo.Session {
	s, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}

	return s
}
