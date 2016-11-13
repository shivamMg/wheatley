package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

const maxMessageSize = 512

// Message is message sent to the user
// type Message struct {
// 	author    string `json:"author"`
// 	message   string `json:"message"`
// 	timestamp string `json:"timestamp"`
// }

// Client is instance of a WS connection at server
type Client struct {
	name string
	send chan []byte
	room *Room
	conn *websocket.Conn
}

// writePump sends data to WS connection
func (c *Client) writePump() {
	for {
		select {
		case strMsg, ok := <-c.send:
			if !ok {
				// Room closed channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(strMsg)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
			// case: Ping
		}
	}
}

// readPump reads data from WS connection and broadcasts in room
func (c *Client) readPump() {
	defer func() {
		c.room.leave <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	for {
		_, strMsg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				fmt.Printf("Error: %v\n", err)
			}
			break
		}
		c.room.broadcast <- strMsg
	}
}
