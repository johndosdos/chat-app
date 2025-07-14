package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/johndosdos/chat_app/internal/routes"
)

func main() {
	const port = "8600"
	const devPort = "5173"
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
	router.Handle("/ws", routes.WebsocketHandler(devPort))

	log.Fatal(server.ListenAndServe())
}
