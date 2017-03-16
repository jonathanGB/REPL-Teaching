package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jonathanGB/REPL-Teaching/app/auth"
	"github.com/jonathanGB/REPL-Teaching/app/models"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"net/http"
	"time"
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
		"userId": user.Id.Hex(),
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
		user.Email,
		fileExtension,
		fileContent,
		readableByteSize(fileSize),
		isPrivateFile == "private",
		time.Now(),
		map[string]bool{},
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

		if (fileOwner == uId) == status { // X-NOR
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
	gId := c.MustGet("group").(*models.GroupInfo).Id

	c.HTML(http.StatusOK, "editor", gin.H{
		"title": fmt.Sprintf("edit %s", file.Name),
		"groupId": gId.Hex(),
		"editor": gin.H{
			"fileName":      file.Name,
			"fileExtension": file.Extension,
			"fileContent":   base64.StdEncoding.EncodeToString(file.Content),
			"privateFile":   file.IsPrivate,
			"isFileOwner":   file.Owner == uId,
		},
	})
}

func (fc *FileController) CloneFile(c *gin.Context) {
	file := c.MustGet("file").(*models.File)
	fileName := c.PostForm("cloneFileName")
	user := c.MustGet("user").(*auth.PublicUser)
	gId := c.MustGet("group").(*models.GroupInfo).Id

	if found, _ := file.ClonedBy[user.Id.Hex()]; found {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "déjà cloné",
		})
		return
	}

	if alreadyFile := fc.model.IsThereUserFile(fileName, gId, user.Id); alreadyFile {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nom de fichier déjà utilisé",
		})
		return
	}

	fc.model.AddCloner(user.Id, gId, file.Id)

	fId := bson.NewObjectId()
	clonedFile := models.File{
		fId,
		fileName,
		user.Id,
		user.Name,
		user.Email,
		file.Extension,
		file.Content,
		file.Size,
		file.IsPrivate,
		time.Now(),
		map[string]bool{},
	}

	if err := fc.model.AddFile(&clonedFile, gId); err != nil {
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

func (fc *FileController) EditorWSHandler(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error connecting websocket %+v", err)
		return
	}

	fmt.Println("connected")
	user := c.MustGet("user").(*auth.PublicUser)
	file := c.MustGet("file").(*models.File)
	gId := c.MustGet("group").(*models.GroupInfo).Id
	wsPayload := struct {
		Type string `json:"type"`
		Content string `json:"content"`
		NewStatus bool `json:"newStatus"`
	}{}
	wsResponse := struct {
		Type string `json:"type"`
		Err bool `json:"err"`
	}{}
	// TODO: Register user to hub

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("closed")
			conn.Close()
			break
		}

		if err := json.Unmarshal(msg, &wsPayload); err != nil {
			conn.WriteMessage(t, []byte("\"bad payload\""))
			continue
		}

		if user.Id != file.Owner && wsPayload.Type != "run" {
			conn.WriteMessage(t, []byte("\"non-owner cannot modify file\""))
			conn.Close()
			break
		}

		switch wsPayload.Type {
		case "run":
			fmt.Println("run") // TODO: link with run service
		case "update-content":
			newContent := []byte(wsPayload.Content)
			newContentSize := len(newContent)

			if newContentSize > MAX_FILE_SIZE {
				err = fmt.Errorf("file too big")
			} else {
				err = fc.model.UpdateFile("content", gId, file.Id, newContent, readableByteSize(int64(newContentSize)))
			}
		case "update-status":
			err = fc.model.UpdateFile("isPrivate", gId, file.Id, wsPayload.NewStatus)
		}

		wsResponse.Type = wsPayload.Type
		wsResponse.Err = (err != nil)

		if err != nil {
			fmt.Println("err in ws", err)
		}
		marshalledResponse, _ := json.Marshal(&wsResponse)
		conn.WriteMessage(t, marshalledResponse)
	}
}

func readableByteSize(size int64) string {
	if size < 1000 {
		return fmt.Sprintf("%dB", size)
	}

	kiloSize := float32(size) / 1000
	return fmt.Sprintf("%.1fkB", kiloSize)
}
