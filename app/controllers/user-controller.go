package controllers

import (
	"github.com/jonathanGB/REPL-Teaching/app/models"
	"github.com/jonathanGB/REPL-Teaching/app/auth"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"os"
)

const MAX_AGE int = 3600 * 24

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
	role := c.PostForm("role")

	if name == "" || email == "" || pwd == "" || role == "" {
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
			role,
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

func (uc *UserController) LoginUser(c *gin.Context) {
	email := c.PostForm("email")
	pwd := c.PostForm("password")

	if email == "" || pwd == "" {
		c.HTML(http.StatusBadRequest, "login", gin.H{
			"title": "login",
			"error": "Paramètre absent",
		})
		return
	}

	user, err := uc.model.FindOne(email, pwd)
	if user.Id == "" || err != nil {
		c.HTML(http.StatusBadRequest, "login", gin.H{
			"title": "login",
			"error": "Mauvaise combinaison",
		})
		return
	}
	token, err := auth.MarshalToken(user.Name, user.Id.Hex(), user.Role)
	if err != nil {
		c.HTML(http.StatusBadRequest, "login", gin.H{
			"title": "login",
			"error": "Erreur sauvegarde session",
		})
		return
	}

	c.SetCookie("auth", token, MAX_AGE, "", "", os.Getenv("env") == "prod", true)
	c.Redirect(http.StatusSeeOther, "/groups/")
}
