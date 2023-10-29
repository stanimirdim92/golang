package main

import (
	"chat/client"
	"net/http"
)

func main() {
	manager := client.NewManager()
	go manager.Start()

	startServer(manager)
}

func startServer(manager *client.Manager) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		manager.WsPage(w, r)
	})

	http.ListenAndServe(":80", nil)
}
