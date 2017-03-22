package controllers

import (
	"github.com/jonathanGB/REPL-Teaching/app/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Message struct {
	fId  bson.ObjectId
	uId  bson.ObjectId
	Data []byte
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
	toEditorObservers chan<- *Message
	toStudentsInMenu  chan<- *Message
	toTeacherInMenu   chan<- *Message

	// register a Client to a specific group
	registerEditorObserver chan<- *Client
	registerStudentInMenu  chan<- *Client
	registerTeacherInMenu  chan<- *Client

	// unregister a Client to a specific group
	unRegisterEditorObserver chan<- *Client
	unRegisterStudentInMenu  chan<- *Client
	unRegisterTeacherInMenu  chan<- *Client
}

func NewHub(s *mgo.Session) Hub {
	h := Hub{}

	gm := models.NewGroupModel(s.Copy())
	gIds := gm.GetAllGroupIds()

	for _, gId := range gIds {
		h[gId.Id] = &Class{
			make(map[bson.ObjectId][]*Client),
			make(map[bson.ObjectId]*Client),
			nil,

			make(chan<- *Message),
			make(chan<- *Message),
			make(chan<- *Message),

			make(chan<- *Client),
			make(chan<- *Client),
			make(chan<- *Client),

			make(chan<- *Client),
			make(chan<- *Client),
			make(chan<- *Client),
		}
		go class.run()
	}
}

func (c *Class) run() {
	for {
		select {
		case client := <-c.registerEditorObserver:
			observers, ok := c.editorObservers[client.fId]
			if !ok {
				observers = make([]*Client)
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
			for i := 0; i < observers; i++ {
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
