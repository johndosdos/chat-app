package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan message)
	mu        = sync.Mutex{}
)

type message struct {
	sender  *websocket.Conn
	content []byte
}

func main() {
	port := ":8080"

	http.Handle("/", http.FileServer(http.Dir("./client")))
	http.HandleFunc("/ws", handleConnection)

	go handleMessages()

	log.Println("Server starting at port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[error] failed to upgrade connection to WebSocket: %v", err)
		return
	}
	defer ws.Close()

	// We need to prevent race conditions when multiple clients are connecting.
	// to our server.
	mu.Lock()
	clients[ws] = true
	mu.Unlock()

	// Process received data for broadcasting to clients.
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			// If the client was disconnected.
			mu.Lock()
			delete(clients, ws)
			mu.Unlock()
			return
		}

		// Send data to the broadcast channel.
		broadcast <- message{
			sender:  ws,
			content: p,
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		// Acquire a lock to prevent race conditions.
		mu.Lock()
		for client := range clients {
			if msg.sender != client {
				err := client.WriteMessage(websocket.TextMessage, msg.content)
				if err != nil {
					log.Printf("[Error] failed to write data to client")
					delete(clients, client)
				}
			}
		}
		mu.Unlock()
	}
}
