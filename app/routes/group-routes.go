package routes

import (
	"github.com/jonathanGB/REPL-Teaching/app/auth"
	"github.com/jonathanGB/REPL-Teaching/app/controllers"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	//"net/http"
)

func GroupRoutes(router *gin.Engine, s *mgo.Session) {
	gc := controllers.NewGroupController(s)

	g := router.Group("/groups", auth.IsAuthentified)
	{
		g.GET("/", gc.FindGroups)
		g.POST("/", auth.IsProf("json"), gc.CreateGroup)
	}
}
