package main

import (
	"chat/client"
	"log"
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

	log.Fatal(http.ListenAndServe(":80", nil))
}
