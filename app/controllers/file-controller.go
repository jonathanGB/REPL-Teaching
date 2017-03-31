package controllers

import (
	"bytes"
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
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var (
	ALLOWED_EXTENSIONS = map[string]bool{
		"go":  true,
		"js":  true,
		"py":  true,
		"rb":  true,
		"exs": true,
		"php": true,
	}

	wsupgrader = websocket.Upgrader{
		ReadBufferSize:  10240,
		WriteBufferSize: 10240,
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
	minimal := c.Query("minimal")

	if minimal == "true" {
		c.HTML(http.StatusUnauthorized, "not-found", gin.H{
			"minimal": true,
		})
		return
	}

	files := fc.model.GetGroupFiles(gInfo.Teacher, gInfo.Id, user.Id, user.Role)
	c.HTML(http.StatusOK, "user-files", gin.H{
		"title":     "Files list",
		"filesMenu": true,
		"role":      user.Role,
		"userId":    user.Id.Hex(),
		"group": gin.H{
			"Id":      gInfo.Id.Hex(),
			"Teacher": gInfo.TeacherName,
			"Files":   files,
		},
	})
}

func (fc *FileController) CreateFile(c *gin.Context) {
	var (
		fileSize    int64
		fileContent = make([]byte, 0)
	)

	// get file, verify size
	mFile, _, err := c.Request.FormFile("fileContent")
	if err == nil {
		fileSize, err = mFile.Seek(0, io.SeekEnd)
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

		fileContent = make([]byte, fileSize)
		if _, err := io.ReadFull(mFile, fileContent); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erreur lors de la lecture du fichier",
			})
			return
		}
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
			"redirect": fmt.Sprintf("/groups/%s/file/%s", gId.Hex(), fId.Hex()),
		})
	}
}

func (fc *FileController) IsFileVisible(c *gin.Context) {
	fIdHex := c.Param("fileId")
	group := c.MustGet("group").(*models.GroupInfo)
	minimal := c.Query("minimal")

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
	query := ""
	if minimal != "" {
		query = "?minimal=true"
	}
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/groups/%s/file/%s", group.Id.Hex(), query))
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
	minimal := c.Query("minimal")

	view := "editor"
	if minimal == "true" {
		view = "minimal-editor"
	}

	c.HTML(http.StatusOK, view, gin.H{
		"title":   fmt.Sprintf("edit %s", file.Name),
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
			"redirect": fmt.Sprintf("/groups/%s/file/%s", gId.Hex(), fId.Hex()),
		})
	}
}

func (fc *FileController) WSEditorOwner(c *gin.Context, class *Class) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error connecting websocket %+v", err)
		return
	}

	fmt.Println("connected")
	user := c.MustGet("user").(*auth.PublicUser)
	file := c.MustGet("file").(*models.File)
	gId := c.MustGet("group").(*models.GroupInfo).Id

	if !file.IsPrivate {
		go class.alertEditorStatus("online", user.Role, file.Id)
	}

	for {
		wsPayload := struct {
			Type           string         `json:"type"`
			Content        string         `json:"content"`
			NewStatus      bool           `json:"newStatus"`
			CursorPosition map[string]int `json:"cursorPosition"`
		}{}
		wsResponse := WSResponse{}

		t, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("closed")
			if !file.IsPrivate {
				class.alertEditorStatus("offline", user.Role, file.Id)
			}
			conn.Close()
			break
		}

		if err := json.Unmarshal(msg, &wsPayload); err != nil {
			conn.WriteMessage(t, []byte("\"bad payload\""))
			continue
		}

		switch wsPayload.Type {
		case "run":
			if wsPayload.Content == "" {
				err = fmt.Errorf("fichier vide")
			} else {
				res, err := runScript(&wsPayload.Content, file.Extension)
				wsResponse.Data = *res
				wsResponse.Err = (err != nil)
			}
		case "update-content":
			newContent := []byte(wsPayload.Content)
			newContentSize := len(newContent)
			newReadableContentSize := readableByteSize(int64(newContentSize))

			if newContentSize > MAX_FILE_SIZE {
				err = fmt.Errorf("file too big")
			} else {
				err = fc.model.UpdateFile("content", gId, file.Id, newContent, newReadableContentSize)

				// alert class
				if err == nil {
					go class.alertContentUpdate(user, file.Id, wsPayload.Content, wsPayload.CursorPosition, file.IsPrivate, newReadableContentSize, time.Now().Format("02 Jan 15:04"))
				}
			}
			wsResponse.Err = (err != nil)
		case "update-status":
			err = fc.model.UpdateFile("isPrivate", gId, file.Id, wsPayload.NewStatus)
			wsResponse.Err = (err != nil)
			wsResponse.Data = wsPayload.NewStatus

			if err == nil {
				go class.alertStatusUpdate(user, file.Id, wsPayload.NewStatus)
				file.IsPrivate = wsPayload.NewStatus
			}
		}

		wsResponse.Type = wsPayload.Type

		if err != nil {
			fmt.Println("err in ws", err)
		}
		marshalledResponse, _ := json.Marshal(&wsResponse)
		conn.WriteMessage(t, marshalledResponse)
	}
}

func (fc *FileController) WSEditorObserver(c *gin.Context, class *Class) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error connecting websocket %+v", err)
		return
	}

	fmt.Println("connected observer")
	user := c.MustGet("user").(*auth.PublicUser)
	file := c.MustGet("file").(*models.File)

	client := Client{
		class,
		file.Id,
		user.Id,
		conn,
		make(chan []byte),
	}

	class.registerEditorObserver <- &client
	go client.writePump()

	for {
		wsPayload := struct {
			Type      string `json:"type"`
			Content   string `json:"content"`
			NewStatus bool   `json:"newStatus"`
		}{}
		wsResponse := WSResponse{}

		t, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("closed")
			class.unRegisterEditorObserver <- &client

			break
		}

		if err := json.Unmarshal(msg, &wsPayload); err != nil {
			conn.WriteMessage(t, []byte("\"bad payload\""))
			continue
		}

		switch wsPayload.Type {
		case "run":
			if wsPayload.Content == "" {
				err = fmt.Errorf("fichier vide")
			} else {
				res, err := runScript(&wsPayload.Content, file.Extension)
				wsResponse.Data = *res
				wsResponse.Err = (err != nil)
			}
		}

		wsResponse.Type = wsPayload.Type

		if err != nil {
			fmt.Println("err in ws", err)
		}
		marshalledResponse, _ := json.Marshal(&wsResponse)
		client.send <- marshalledResponse
	}
}

func (fc *FileController) WSInMenu(c *gin.Context, class *Class) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("Error connecting websocket %+v\n", err)
		return
	}

	fmt.Println("connected in-menu")
	user := c.MustGet("user").(*auth.PublicUser)

	client := Client{
		class,
		"",
		user.Id,
		conn,
		make(chan []byte),
	}
	editorsQuery := EditorQuery{
		user.Role,
		make([]string, 0),
		make(chan bool, 1),
	}

	if user.Role == "teacher" {
		class.registerTeacherInMenu <- &client
	} else {
		class.registerStudentInMenu <- &client
	}

	class.getPublicEditors <- &editorsQuery
	<-editorsQuery.done

	liveResponse := WSResponse{
		"live-editing",
		false,
		map[string]interface{}{
			"files":  editorsQuery.editors,
			"status": "online",
		},
	}
	fmt.Println("after...", liveResponse)
	marshalledResponse, _ := json.Marshal(&liveResponse)

	m := Message{
		"",
		user.Id,
		marshalledResponse,
	}

	if user.Role == "teacher" {
		class.toTeacherInMenu <- &m
	} else {
		class.toStudentsInMenu <- &m
	}

	client.writePump()
}

func readableByteSize(size int64) string {
	if size < 1000 {
		return fmt.Sprintf("%dB", size)
	}

	kiloSize := float32(size) / 1000
	return fmt.Sprintf("%.1fkB", kiloSize)
}

func runScript(script *string, extension string) (*string, error) {
	apiUrl := "http://localhost:8081/run"
	data := url.Values{}
	data.Set("extension", extension)
	data.Set("content", *script)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", apiUrl, bytes.NewBufferString(data.Encode()))
	r.Header.Add("content-type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	strBody := string(body)

	if resp.StatusCode != http.StatusOK {
		return &strBody, fmt.Errorf("erreur provenant du service de run")
	} else {
		return &strBody, nil
	}
}
