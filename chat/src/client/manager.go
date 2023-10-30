package client

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type Manager struct {
	// Registered clients.
	clients map[*Client]bool
	// Inbound messages from the clients.
	broadcast chan []byte
	// Register requests from the clients.
	register chan *Client
	// Unregister requests from clients.
	unregister chan *Client
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		//fmt.Print(r.Header, r.Host, r.RemoteAddr, r.RequestURI)
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (manager *Manager) Start() *Manager {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true // put the Client inside a map
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

func (manager *Manager) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}

func (manager *Manager) WsPage(res http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(res, req, nil)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	//defer conn.Close()

	c := NewClient(conn)

	manager.register <- c

	go manager.Read(c)
	go manager.Write(c)
}

func (manager *Manager) Read(client *Client) {
	defer func() {
		manager.unregister <- client
	}()

	for {
		_, message, err := client.socket.ReadMessage()
		if err != nil {
			manager.unregister <- client
			break
		}
		jsonMessage, _ := json.Marshal(&Message{
			Id:        client.RandomString(32),
			Sender:    client.id,
			Content:   string(message),
			Timestamp: time.Now().Unix(),
		})
		manager.broadcast <- jsonMessage
	}
}

func (manager *Manager) Write(client *Client) {
	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				client.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			client.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
