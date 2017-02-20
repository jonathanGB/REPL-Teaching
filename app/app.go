package main

import (
	"fmt"
	"github.com/gin-contrib/multitemplate"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
)

var s = initDB()

func templateRender() multitemplate.Render {
	r := multitemplate.New()

	r.AddFromFiles("signup", "templates/layout.gohtml", "templates/signup.gohtml")
	r.AddFromFiles("login", "templates/layout.gohtml", "templates/login.gohtml")
	r.AddFromFiles("logout", "templates/layout.gohtml", "templates/logout.gohtml")

	return r
}

func getMainEngine() *gin.Engine {
	router := gin.Default()
	router.Static("/public", "./public")
	router.StaticFile("/favicon.ico", "./public/img/favicon.ico")
	router.HTMLRender = templateRender()

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup", gin.H{
			"title": "signup",
		})
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login", gin.H{
			"title": "login",
		})
	})

	router.GET("/logout", func(c *gin.Context) {
		c.HTML(http.StatusOK, "logout", gin.H{
			"title": "logout",
			"name":  "user1234",
		})
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
	defer s.Close() // mongodb session

	getMainEngine().Run(":8080")
}
