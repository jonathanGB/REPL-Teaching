package controllers

import (
	"fmt"
	"github.com/jonathanGB/REPL-Teaching/app/auth"
	"github.com/jonathanGB/REPL-Teaching/app/models"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"net/http"
)

var (
	ALLOWED_EXTENSIONS = map[string]bool{
		"go": true,
		"js": true,
	}
)

const (
	MAX_FILE_SIZE = 10000 // 10kB
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

// TODO: complete getting files from db
func (fc *FileController) ShowGroupFiles(c *gin.Context) {
	gInfo := c.MustGet("group").(*models.GroupInfo)

	c.HTML(http.StatusOK, "user-files", gin.H{
		"title": "Files list",
		"group": gin.H{
			"Id": gInfo.Id.Hex(),
		},
	})
}

func (fc *FileController) CreateFile(c *gin.Context) {
	// get file, verify size
	mFile, _, err := c.Request.FormFile("fileContent")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file is absent / could not be parsed",
		})
		return
	}

	fileSize, err := mFile.Seek(0, io.SeekEnd)
	mFile.Seek(0, io.SeekStart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal error",
		})
		return
	}

	if fileSize > MAX_FILE_SIZE {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "size limit is 10kB",
		})
		return
	}

	// get other params
	fileName := c.PostForm("fileName")
	fileExtension := c.PostForm("fileExtension")
	isPrivateFile := c.PostForm("isPrivate")

	if fileName == "" || fileExtension == "" || !ALLOWED_EXTENSIONS[fileExtension] || (isPrivateFile != "public" && isPrivateFile != "private") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "paramètre(s) manquant(s)",
		})
		return
	}

	gId := c.MustGet("group").(*models.GroupInfo).Id
	uId := c.MustGet("user").(*auth.PublicUser).Id

	if alreadyFile := fc.model.IsThereUserFile(fileName, gId, uId); alreadyFile {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nom de fichier déjà utilisé",
		})
		return
	}

	fId := bson.NewObjectId()
	fileContent := make([]byte, fileSize)
	if _, err := io.ReadFull(mFile, fileContent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la lecture du fichier",
		})
		return
	}

	file := models.File{
		fId,
		fileName,
		fId,
		uId,
		fileExtension,
		fileContent,
		isPrivateFile == "private",
	}

	if err := fc.model.AddFile(&file, gId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la création du fichier",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"error":    nil,
			"redirect": fmt.Sprintf("/groups/%s/files/%s", gId.Hex(), fId.Hex()),
		})
	}

}
