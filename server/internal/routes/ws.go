package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

func WebsocketHandler(port string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow the vite dev server for cross-origin resource sharing
		// Enable local network testing

		log.Printf("[server] Establishing WebSocket server, [origin] %v\n", r.Header.Get("Origin"))

		// This handler demonstrates how to safely accept cross origin WebSockets
		// from the origin example.com.
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			OriginPatterns: []string{
				"localhost:" + port,
				"0.0.0.0:" + port,
			},
			// InsecureSkipVerify: true,
		})
		if err != nil {
			log.Printf("[error] %v\n", err)
			return
		}
		log.Println("[server] Client connected")
		defer conn.CloseNow()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for {
			_, data, err := conn.Read(ctx)
			if err != nil {
				if websocket.CloseStatus(err) != -1 {
					log.Printf("[info] Client disconnected: %v\n", err)
				} else {
					log.Printf("[error] Failed to read connection: %v\n", err)
				}
				conn.Close(websocket.StatusNormalClosure, "[info] Connection closed")
				return
			}

			log.Printf("[client msg] %v\n", string(data))
		}
	})
}
