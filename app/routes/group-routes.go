package routes

import (
	"fmt"
	"github.com/jonathanGB/REPL-Teaching/app/auth"
	"github.com/jonathanGB/REPL-Teaching/app/controllers"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"net/http"
)

func GroupRoutes(router *gin.Engine, s *mgo.Session) {
	gc := controllers.NewGroupController(s)

	groups := router.Group("/groups", auth.IsAuthentified)
	{
		groups.GET("/", gc.FindGroups)
		groups.POST("/", auth.IsProf(true, "json"), gc.CreateGroup)

		group := groups.Group("/:groupId", gc.IsGroup)
		{
			group.GET("/", func(c *gin.Context) {
				gId := c.Param("groupId")
				c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/groups/%s/files", gId))
			})

			group.GET("/join", auth.IsProf(false, "html"), gc.IsGroupMember(false), gc.ShowJoiningGroup)
			group.POST("/join", auth.IsProf(false, "json"), gc.IsGroupMember(false), gc.JoinGroup)

			// TODO: dummy response for now
			group.GET("/files", gc.IsGroupMember(true), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"data": c.Param("groupId"),
				})
			})
		}
	}
}
