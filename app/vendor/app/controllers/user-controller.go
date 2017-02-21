package controllers

import (
	"app/models"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type UserController struct {
	model *models.UserModel
}

func NewUserController(s *mgo.Session) *UserController {
	return &UserController{
		models.NewUserModel(s.Copy()),
	}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	pwd := c.PostForm("password")

	if name == "" || email == "" || pwd == "" {
		c.HTML(http.StatusBadRequest, "signup", gin.H{
			"title": "Sign up",
			"error": "Paramètre absent",
		})
		return
	}

	// check for email duplicate, then add to db
	if alreadyIn := uc.model.IsThere(email); !alreadyIn {
		hashedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "signup", gin.H{
				"title": "Sign up",
				"error": "Erreur lors de la création du compte",
				"name":  name,
				"email": email,
			})
			return
		}

		user := models.User{
			bson.NewObjectId(),
			name,
			email,
			hashedPwd,
		}

		if err := uc.model.AddUser(&user); err == nil {
			// show success message
			c.HTML(http.StatusOK, "signedup", gin.H{
				"title": "Signed up",
				"name":  name,
			})
		} else {
			c.HTML(http.StatusInternalServerError, "signup", gin.H{
				"title": "Sign up",
				"error": "Erreur lors de la création du compte",
				"name":  name,
				"email": email,
			})
		}
	} else {
		c.HTML(http.StatusBadRequest, "signup", gin.H{
			"title": "Sign up",
			"error": "Email déjà utilisé",
			"name":  name,
		})
	}
}
