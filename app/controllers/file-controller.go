package controllers

import (
	//"github.com/jonathanGB/REPL-Teaching/app/auth"
	"github.com/jonathanGB/REPL-Teaching/app/models"
	//"golang.org/x/crypto/bcrypt"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"net/http"
)

// TODO: add reference to pool of containers, or have it as a separate controller?
type FileController struct {
	model *models.FileModel
}

func NewFileController(s *mgo.Session) *FileController {
	return &FileController{
		models.NewFileModel(s.Copy()),
	}
}

func (fc *FileController) CreateFile(c *gin.Context) {
	// get multi-part file
	// verify extension & size
	// get isPrivate field
	// verify no collision with user-visible file names
	c.JSON(http.StatusOK, gin.H{
		"data": "to fill",
	})

}
