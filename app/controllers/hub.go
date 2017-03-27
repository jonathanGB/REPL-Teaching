package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/jonathanGB/REPL-Teaching/app/auth"
	"github.com/jonathanGB/REPL-Teaching/app/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Message struct {
	fId  bson.ObjectId
	uId  bson.ObjectId
	Data []byte
}
type WSResponse struct {
	Type string      `json:"type"`
	Err  bool        `json:"err"`
	Data interface{} `json:"data"`
}
type Hub map[bson.ObjectId]*Class // key is groupId
type Class struct {
	// every connection Client must be in one AND ONLY one of these 3 categories
	// due to the iframe preview, any user can have 2 simultaneous ws connections: 1 in the editor and 1 in the menu
	// editor and menu ws connections send/receive different data
	// ----
	//student updates file -> toEditorObservers[fId]... & studentsInMenu[uId]
	//teacher updates file -> toEditorObservers[fId]... & teacherInMenu
	//student adds file -> teacherInMenu
	//teacher adds file -> studentsInMenu...
	editorObservers map[bson.ObjectId][]*Client // all Clients observing a file (excluding owner)
	studentsInMenu  map[bson.ObjectId]*Client   // students in the files menu (maps uId to student in menu)
	teacherInMenu   *Client                     // teacher in the menu

	// send data to a specific class of Clients
	toEditorObservers chan *Message
	toStudentsInMenu  chan *Message
	toTeacherInMenu   chan *Message

	// register a Client to a specific group
	registerEditorObserver chan *Client
	registerStudentInMenu  chan *Client
	registerTeacherInMenu  chan *Client

	// unregister a Client to a specific group
	unRegisterEditorObserver chan *Client
	unRegisterStudentInMenu  chan *Client
	unRegisterTeacherInMenu  chan *Client
}

func NewHub(s *mgo.Session) Hub {
	h := Hub{}

	gm := models.NewGroupModel(s.Copy())
	gIds := gm.GetAllGroupIds()

	for _, gId := range gIds {
		h.registerClass(gId.Id)
	}

	return h
}

func (h Hub) registerClass(gId bson.ObjectId) {
	class := &Class{
		make(map[bson.ObjectId][]*Client),
		make(map[bson.ObjectId]*Client),
		nil,

		make(chan *Message),
		make(chan *Message),
		make(chan *Message),

		make(chan *Client),
		make(chan *Client),
		make(chan *Client),

		make(chan *Client),
		make(chan *Client),
		make(chan *Client),
	}
	h[gId] = class

	go class.run()
}

func (c *Class) alertStatusUpdate(user *auth.PublicUser, fId bson.ObjectId, newStatus bool) {
	res := WSResponse{
		"update-status",
		false,
		newStatus,
	}
	jsonRes, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Error decoding update-status json")
		return
	}

	m := Message{
		fId,
		user.Id,
		jsonRes,
	}

	if user.Role == "teacher" { // alert all students
		m.uId = ""
	}

	c.toEditorObservers <- &m
	c.toTeacherInMenu <- &m
	c.toStudentsInMenu <- &m
}

func (c *Class) alertContentUpdate(user *auth.PublicUser, fId bson.ObjectId, newContent []byte, isPrivate bool, readableSize, lastModified string) {
	// encode file update data (meta & all data separately)
	metaData := make(map[string]interface{})
	metaData["size"] = readableSize
	metaData["lastModified"] = lastModified

	allData := make(map[string]interface{})
	allData["size"] = readableSize
	allData["lastModified"] = lastModified
	allData["content"] = newContent

	// encode ws response
	metaDataRes := WSResponse{
		"update-content",
		false,
		metaData,
	}
	jsonMetaData, err := json.Marshal(metaDataRes)
	if err != nil {
		fmt.Println("Error decoding update-content json")
		return
	}

	allDataRes := WSResponse{
		"update-content",
		false,
		allData,
	}
	jsonAllData, err := json.Marshal(allDataRes)
	if err != nil {
		fmt.Println("Error decoding update-content json")
		return
	}

	// prepare messages to send to channels
	lightM := Message{
		fId,
		user.Id,
		jsonMetaData,
	}
	allM := Message{
		fId,
		user.Id,
		jsonAllData,
	}

	if user.Role == "teacher" {
		lightM.uId = ""
		allM.uId = ""

		c.toTeacherInMenu <- &lightM
	} else {
		c.toStudentsInMenu <- &lightM
	}

	if !isPrivate {
		c.toEditorObservers <- &allM

		// share metadata to students if teacher, and to teacher if student
		if user.Role == "teacher" {
			c.toStudentsInMenu <- &lightM
		} else {
			c.toTeacherInMenu <- &lightM
		}
	}
}

func (c *Class) run() {
	for {
		select {
		case client := <-c.registerEditorObserver:
			observers, ok := c.editorObservers[client.fId]
			if !ok {
				observers = []*Client{}
				c.editorObservers[client.fId] = observers
			}
			c.editorObservers[client.fId] = append(observers, client)

		case client := <-c.registerStudentInMenu:
			c.studentsInMenu[client.uId] = client

		case client := <-c.registerTeacherInMenu:
			c.teacherInMenu = client

		case client := <-c.unRegisterEditorObserver:
			observers, _ := c.editorObservers[client.fId]
			for i, observer := range observers {
				if client == observer {
					c.editorObservers[client.fId] = append(observers[:i], observers[i+1:]...)
					break
				}
			}
			if len(c.editorObservers[client.fId]) == 0 {
				delete(c.editorObservers, client.fId)
			}
			close(client.send)

		case client := <-c.unRegisterStudentInMenu:
			delete(c.studentsInMenu, client.uId)
			close(client.send)

		case client := <-c.unRegisterTeacherInMenu:
			c.teacherInMenu = nil
			close(client.send)

		case message := <-c.toEditorObservers:
			observers, ok := c.editorObservers[message.fId]
			if !ok {
				continue
			}
			for i := 0; i < len(observers); i++ {
				client := observers[i]
				select {
				case client.send <- message.Data:
				default:
					close(client.send)
					observers = append(observers[:i], observers[i+1:]...)
					i--
				}
			}
			c.editorObservers[message.fId] = observers

		case message := <-c.toStudentsInMenu:
			if message.uId == "" { // teacher to all students
				for key, client := range c.studentsInMenu {
					select {
					case client.send <- message.Data:
					default:
						close(client.send)
						delete(c.studentsInMenu, key)
					}
				}
				continue
			}
			client, ok := c.studentsInMenu[message.uId]
			if !ok {
				continue
			}
			select {
			case client.send <- message.Data:
			default:
				close(client.send)
				delete(c.studentsInMenu, message.uId)
			}

		case message := <-c.toTeacherInMenu:
			if c.teacherInMenu == nil {
				continue
			}
			select {
			case c.teacherInMenu.send <- message.Data:
			default:
				close(c.teacherInMenu.send)
				c.teacherInMenu = nil
			}
		}
	}
}
