package client

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type client struct {
	id string
	// The websocket connection.
	socket *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}

func (c client) RandomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

func newClient(conn *websocket.Conn) *client {
	return &client{
		id:     uuid.NewV4().String(),
		socket: conn,
		send:   make(chan []byte),
	}
}
