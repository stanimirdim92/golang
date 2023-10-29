package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
)

type Message struct {
	Sender    *Client `json:"sender,omitempty"`
	Recipient *Client `json:"recipient,omitempty"`
	Content   string  `json:"content,omitempty"`
}

type Client struct {
	id      string
	manager *ClientManager
	// The websocket connection.
	socket *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		fmt.Print(r.Header, r.Host, r.RemoteAddr, r.RequestURI)
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func unregister(c *Client) {
	c.manager.unregister <- c
	c.socket.Close()
}

func (c *Client) read() {
	defer func() {
		unregister(c)
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			unregister(c)
			break
		}
		jsonMessage, _ := json.Marshal(&Message{Sender: c, Content: string(message)})
		c.manager.broadcast <- jsonMessage
	}
}

func (c *Client) write() {
	defer c.socket.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func wsPage(res http.ResponseWriter, req *http.Request, manager *ClientManager) {
	conn, err := upgrader.Upgrade(res, req, nil)

	if err != nil {
		log.Println(err)
		http.NotFound(res, req)
		return
	}

	client := &Client{
		id:      uuid.NewV4().String(),
		socket:  conn,
		send:    make(chan []byte),
		manager: manager,
	}

	manager.register <- client

	go client.read()
	go client.write()
}
