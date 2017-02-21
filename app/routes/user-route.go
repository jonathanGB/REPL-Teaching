package route

import (
	"github.com/jonathanGB/REPL-Teaching/app/controllers"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"net/http"
)

func UserRoutes(router *gin.Engine, s *mgo.Session) {
	uc := controllers.NewUserController(s)

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/users")
	})

	router.GET("/users", func(c *gin.Context) {
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

	router.POST("/users", uc.CreateUser)
}
