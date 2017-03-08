package controllers

import (
	"encoding/base64"
	"fmt"
	"github.com/gorilla/websocket"
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

	wsupgrader = websocket.Upgrader{
		ReadBufferSize:  10240,
		WriteBufferSize: 1024,
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
	user := c.MustGet("user").(*auth.PublicUser)

	files := fc.model.GetGroupFiles(gInfo.Teacher, gInfo.Id, user.Id, user.Role)

	c.HTML(http.StatusOK, "user-files", gin.H{
		"title": "Files list",
		"role":  user.Role,
		"group": gin.H{
			"Id":      gInfo.Id.Hex(),
			"Teacher": gInfo.TeacherName,
			"Files":   files,
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
	user := c.MustGet("user").(*auth.PublicUser)

	if alreadyFile := fc.model.IsThereUserFile(fileName, gId, user.Id); alreadyFile {
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
		user.Id,
		user.Name,
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

func (fc *FileController) IsFileVisible(c *gin.Context) {
	fIdHex := c.Param("fileId")
	group := c.MustGet("group").(*models.GroupInfo)

	if !bson.IsObjectIdHex(fIdHex) {
		c.Abort()
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/groups/%s/files", group.Id.Hex()))
		return
	}
	fId := bson.ObjectIdHex(fIdHex)

	file, err := fc.model.GetFile(fId, group.Id)
	if err != nil {
		c.Abort()
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/groups/%s/files", group.Id.Hex()))
		return
	}

	user := c.MustGet("user").(*auth.PublicUser)
	// teacher can see all public files of the group and his own
	// student can see own files, or public by teacher
	if file.Owner == user.Id || (user.Role == "teacher" || file.Owner == group.Teacher) && !file.IsPrivate {
		c.Set("file", file)
		c.Next()
		return
	}

	c.Abort()
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/groups/%s/files", group.Id.Hex()))
}

func (fc *FileController) IsFileOwner(status bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileOwner := c.MustGet("file").(*models.File).Owner
		uId := c.MustGet("user").(*auth.PublicUser).Id
		gId := c.MustGet("group").(*models.GroupInfo).Id

		if fileOwner == uId && status || fileOwner != uId && !status {
			c.Next()
			return
		}

		c.Abort()
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/groups/%s/files", gId.Hex()))
	}
}

func (fc *FileController) ShowFile(c *gin.Context) {
	file := c.MustGet("file").(*models.File)
	uId := c.MustGet("user").(*auth.PublicUser).Id

	c.HTML(http.StatusOK, "editor", gin.H{
		"title": fmt.Sprintf("edit %s", file.Name),
		"editor": gin.H{
			"fileName":      file.Name,
			"fileExtension": file.Extension,
			"fileContent":   base64.StdEncoding.EncodeToString(file.Content),
			"privateFile":   file.IsPrivate,
			"isFileOwner":   file.Owner == uId,
		},
	})
}

func (fc *FileController) EditorWSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error connecting websocket %+v", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("closed")
			conn.Close()
			break
		}
		fmt.Println(t)

		conn.WriteMessage(t, msg)
	}
}
