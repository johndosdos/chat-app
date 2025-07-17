package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/johndosdos/chat_app/internal/routes"
)

func main() {
	const port = "8080"
	const vitePort = "5173"

	router := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Handle root endpoint.
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received for %s from host %s", r.URL.Path, r.Host)
		fmt.Fprint(w, "Hello from server!\n")
	})

	// Handle websocket endpoint.
	router.Handle("/ws", routes.WebsocketHandler(vitePort))

	log.Printf("Starting server on %v...\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}
