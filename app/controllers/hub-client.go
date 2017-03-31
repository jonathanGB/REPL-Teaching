package controllers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2/bson"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  10240,
	WriteBufferSize: 10240,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	class *Class

	fId, uId bson.ObjectId

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	for {
		message, ok := <-c.send
		fmt.Println("new messsage incoming")

		if !ok {
			// The hub closed the channel.
			fmt.Println("error in message")
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		c.conn.WriteMessage(websocket.TextMessage, message)
		fmt.Println("sent message!")
	}
}
