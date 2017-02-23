package controllers

import (
	"github.com/jonathanGB/REPL-Teaching/app/auth"
	"github.com/jonathanGB/REPL-Teaching/app/models"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"net/http"
)

type GroupController struct {
	model *models.GroupModel
}

func NewGroupController(s *mgo.Session) *GroupController {
	return &GroupController{
		models.NewGroupModel(s.Copy()),
	}
}

func (gc *GroupController) FindGroups(c *gin.Context) {
	user, _ := c.Get("user")

	groups := gc.model.GetUserGroups(user.(*auth.PublicUser).Id)

	c.HTML(http.StatusOK, "user-groups", gin.H{
		"title": "groups dashboard",
		"user": user.(*auth.PublicUser),
		"data": groups,
	})
}
