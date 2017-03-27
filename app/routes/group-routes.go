package routes

import (
	"fmt"
	"github.com/jonathanGB/REPL-Teaching/app/auth"
	"github.com/jonathanGB/REPL-Teaching/app/controllers"
	"github.com/jonathanGB/REPL-Teaching/app/models"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"net/http"
)

func GroupRoutes(router *gin.Engine, s *mgo.Session, hub controllers.Hub) {
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
				files.GET("/", gc.IsGroupMember(true), fc.ShowGroupFiles)
				files.POST("/", gc.IsGroupMember(true), fc.CreateFile)

				// TODO: add ws handlers in the menu

				file := files.Group("/:fileId", fc.IsFileVisible)
				{
					file.GET("/", fc.ShowFile)

					file.GET("/ws", func(c *gin.Context) {
						uId := c.MustGet("user").(*auth.PublicUser).Id
						fOwner := c.MustGet("file").(*models.File).Owner
						gId := c.MustGet("group").(*models.GroupInfo).Id

						if uId == fOwner {
							fc.WSEditorOwner(c, hub[gId])
						} else {
							fc.WSEditorObserver(c, hub[gId])
						}
					})

					file.POST("/clone", auth.IsProf(false, "json"), fc.IsFileOwner(false), fc.CloneFile)
				}
			}
		}
	}
}
