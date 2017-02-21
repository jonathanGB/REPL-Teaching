package route

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

func GroupRoutes(router *gin.Engine) {
	g := router.Group("/groups")
	{
		g.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "in construction")
		})
	}
}
