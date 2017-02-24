package controllers

import (
	"github.com/jonathanGB/REPL-Teaching/app/auth"
	"github.com/jonathanGB/REPL-Teaching/app/models"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	user := c.MustGet("user").(*auth.PublicUser)

	groups := gc.model.GetUserGroups(user.Id)

	c.HTML(http.StatusOK, "user-groups", gin.H{
		"title": "groups dashboard",
		"user":  user,
		"data":  groups,
	})
}

func (gc *GroupController) CreateGroup(c *gin.Context) {
	user := c.MustGet("user").(*auth.PublicUser)
	gName := c.PostForm("groupName")
	gPwd := c.PostForm("groupPassword")

	if alreadyGroup := gc.model.IsThereGroup(gName, user.Id); alreadyGroup {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nom déjà utilisé",
		})
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(gPwd), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la création du groupe",
		})
		return
	}

	gId := bson.NewObjectId()
	group := models.Group{
		gId,
		gName,
		user.Id,
		user.Name,
		[]models.File{},
		hashedPwd,
	}

	if err := gc.model.AddGroup(&group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la création du groupe",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"error": nil,
			"data":  gId,
		})
	}
}
