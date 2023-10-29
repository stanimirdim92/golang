package client

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type Manager struct {
	// Registered clients.
	clients map[*client]bool
	// Inbound messages from the clients.
	broadcast chan []byte
	// Register requests from the clients.
	register chan *client
	// Unregister requests from clients.
	unregister chan *client
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
		clients:    make(map[*client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
	}
}

func (manager *Manager) Start() *Manager {
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

func (manager *Manager) send(message []byte, ignore *client) {
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

	client := newClient(conn)

	manager.register <- client

	go manager.Read(client)
	go manager.Write(client)
}

func (manager *Manager) Read(client *client) {
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

func (manager *Manager) Write(client *client) {
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
