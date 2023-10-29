package main

import "net/http"

func main() {
	manager := NewManager()
	go manager.start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsPage(w, r, manager)
	})

	http.ListenAndServe(":12345", nil)
}
