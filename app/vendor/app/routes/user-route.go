package route

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

func UserRoutes(router *gin.Engine) {
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
}
