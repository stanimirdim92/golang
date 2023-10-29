package main

import (
	"encoding/json"
)

type ClientManager struct {
	// Registered clients.
	clients map[*Client]bool
	// Inbound messages from the clients.
	broadcast chan []byte
	// Register requests from the clients.
	register chan *Client
	// Unregister requests from clients.
	unregister chan *Client
}

func NewManager() *ClientManager {
	return &ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}
func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true // put the client inside a map
			jsonMessage, _ := json.Marshal(&Message{Content: "Client has connected."})
			manager.send(jsonMessage, conn)

		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				delete(manager.clients, conn)
				close(conn.send)
				jsonMessage, _ := json.Marshal(&Message{Content: "Client has disconnected."})
				manager.send(jsonMessage, conn)
			}

		case message := <-manager.broadcast:
			for conn := range manager.clients {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}
