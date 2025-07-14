package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/johndosdos/chat_app/internal/routes"
)

func main() {
	const port = "8080"
	fmt.Printf("Starting server on port %v...\n", port)

	router := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Handle root endpoint.
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello from server!\n")
	})

	// Handle websocket endpoint.
	router.HandleFunc("/ws", routes.WebsocketHandler)

	log.Fatal(server.ListenAndServe())
}
