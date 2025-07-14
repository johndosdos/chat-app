package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

func WebsocketHandler(port string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow the vite dev server for cross-origin resource sharing
		devUrl := "localhost:" + port
		fmt.Println("[server] Establishing WebSocket server")

		// This handler demonstrates how to safely accept cross origin WebSockets
		// from the origin example.com.

		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			OriginPatterns: []string{devUrl},
		})
		if err != nil {
			log.Printf("[error] %v\n", err)
			return
		}
		defer conn.Close(websocket.StatusNormalClosure, "Connection closed")
	})
}
