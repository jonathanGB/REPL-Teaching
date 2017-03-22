package routes

import (
	"fmt"
	"github.com/jonathanGB/REPL-Teaching/app/auth"
	"github.com/jonathanGB/REPL-Teaching/app/controllers"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"net/http"
)

func GroupRoutes(router *gin.Engine, s *mgo.Session, hub *controllers.Hub) {
	gc := controllers.NewGroupController(s)
	fc := controllers.NewFileController(s)

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

			files := group.Group("/files", gc.IsGroupMember(true))
			{
				// TODO: dummy response for now
				files.GET("/", gc.IsGroupMember(true), fc.ShowGroupFiles)
				files.POST("/", gc.IsGroupMember(true), fc.CreateFile)

				file := files.Group("/:fileId", fc.IsFileVisible)
				{
					file.GET("/", fc.ShowFile)

					file.GET("/ws", func(c *gin.Context) {
						fc.EditorWSHandler(c, hub)
					})
					// file.PUT("/", fc.IsFileOwner(true), fc.UpdateFile)

					file.POST("/clone", auth.IsProf(false, "json"), fc.IsFileOwner(false), fc.CloneFile)
				}
			}
		}
	}
}
