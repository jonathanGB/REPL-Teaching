package main

import (
	"fmt"
	"github.com/gin-contrib/multitemplate"
	"github.com/jonathanGB/REPL-Teaching/app/routes"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
)

var s = initDB()

func templateRender() multitemplate.Render {
	r := multitemplate.New()

	r.AddFromFiles("signup", "templates/layout.gohtml", "templates/signup.gohtml")
	r.AddFromFiles("login", "templates/layout.gohtml", "templates/login.gohtml")
	r.AddFromFiles("logout", "templates/layout.gohtml", "templates/logout.gohtml")
	r.AddFromFiles("signedup", "templates/layout.gohtml", "templates/signedup.gohtml")
	r.AddFromFiles("user-groups", "templates/layout.gohtml", "templates/user-groups.gohtml")
	r.AddFromFiles("join-group", "templates/layout.gohtml", "templates/join-group.gohtml")

	return r
}

func getMainEngine() *gin.Engine {
	app := gin.Default()
	app.Static("/public", "./public")
	app.StaticFile("/favicon.ico", "./public/img/favicon.ico")
	app.HTMLRender = templateRender()

	// add routes
	routes.FooBarRoutes(app, s)
	routes.UserRoutes(app, s)
	routes.GroupRoutes(app, s)

	fmt.Println("\n") // empty buffer in output
	return app
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
