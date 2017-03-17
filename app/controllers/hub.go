package controllers

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/jonathanGB/REPL-Teaching/app/models"
)

type Hub map[bson.ObjectId]*Class // key is groupId
type Class struct {
	students map[*Client]bool
	teacher *Client
	toStudents chan<- []byte
	toTeacher chan<- []byte
	register chan<- *Client
	unregister chan<- *Client
}

func NewHub(s *mgo.Session) Hub {
	h := Hub{}

	gm := models.NewGroupModel(s.Copy())
	gIds := gm.GetAllGroupIds()

	for _, gId := range gIds {
		class := &Class{
			make(map[*Client]bool),
			nil,
			make(chan<- []byte),
			make(chan<- []byte),
			make(chan<- *Client),
			make(chan<- *Client),
		}
		h[gId.Id] = class
		go class.run()
	}
}

func (c *Class) run() {
	for {
		select {
		case client := <-c.register:
			if client.Type == "teacher" {
				c.teacher = client
			} else {
				c.students[client] = true
			}
		case client := <-h.unregister:
			if client == c.teacher {
				c.teacher = nil
				close(client.send)
			} else if _, ok := c.students[client]; ok {
				delete(c.clients, client)
				close(client.send)
			} else {
				close(client.send)
			}
		case message := <-c.toTeacher && c.teacher != nil:
			select {
			case c.teacher.send <- message:
			default:
				close(c.teacher.send)
				c.teacher = nil
		case message := <-c.toStudents:
			for student := range c.students {
				select {
				case student.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
