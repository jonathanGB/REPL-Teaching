package routes

import (
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
)

func FooBarRoutes(router *gin.Engine, s *mgo.Session) {
	router.GET("/hello", func(c *gin.Context) {
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
}
