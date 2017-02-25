package controllers

import (
	"fmt"
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

	if err := gc.model.AddGroup(&group, user.Id); err != nil {
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

func (gc *GroupController) IsGroup(c *gin.Context) {
	gIdHex := c.Param("groupId")
	fmt.Println(gIdHex)
	if !bson.IsObjectIdHex(gIdHex) {
		c.Abort()
		c.Redirect(http.StatusSeeOther, "/groups")
		return
	}
	gId := bson.ObjectIdHex(gIdHex)

	// check if in db, store rendered group data in context
	gInfo := gc.model.GetGroupInfo(gId)
	if gInfo.Id == "" {
		c.Abort()
		c.Redirect(http.StatusSeeOther, "/groups")
		return
	}

	c.Set("group", gInfo)
	c.Next()
}

func (gc *GroupController) IsGroupMember(status bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*auth.PublicUser)
		gId := c.Param("groupId")
		found := false

		// check if gId in group slice of user
		userGroups := gc.model.GetUserGroups(user.Id)
		for _, userGroup := range userGroups {
			if userGroup.Id == gId {
				found = true
				break
			}
		}

		if found == status {
			c.Next()
			return
		}

		c.Abort()
		if status { // not member, but should
			c.Redirect(http.StatusSeeOther, fmt.Sprintf("/groups/%s/join", gId))
		} else { // member, but shouldn't
			c.Redirect(http.StatusSeeOther, fmt.Sprintf("/groups/%s/files", gId))
		}
	}
}

func (gc *GroupController) ShowJoiningGroup(c *gin.Context) {
	gInfo := c.MustGet("group").(*models.GroupInfo)

	c.HTML(http.StatusOK, "join-group", gin.H{
		"title": fmt.Sprintf("Join %s", gInfo.Name),
		"group": gin.H{
			"Id":   gInfo.Id.Hex(),
			"Name": gInfo.Name,
		},
	})
}

func (gc *GroupController) JoinGroup(c *gin.Context) {
	userId := c.MustGet("user").(*auth.PublicUser).Id
	gInfo := c.MustGet("group").(*models.GroupInfo)
	pwd := c.PostForm("groupPassword")

	if err := bcrypt.CompareHashAndPassword(gInfo.Password, []byte(pwd)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Mauvais mot de passe",
		})
		return
	}

	if err := gc.model.JoinGroup(gInfo.Id, userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur en joignant le groupe",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":    nil,
		"redirect": fmt.Sprintf("/groups/%s/files", gInfo.Id.Hex()),
	})
}
