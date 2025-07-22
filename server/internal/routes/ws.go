package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/coder/websocket"
	"github.com/johndosdos/chat_app/internal/clients"
)

func WebsocketHandler(port string, cls *clients.Clients) http.Handler {
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

		ctx := context.Background()
		serverMsg := "hello clients, from server"

		cl := clients.NewClient(conn)
		cls.Add(cl.Id, cl)

		defer cl.Conn.Close(websocket.StatusNormalClosure, "Connection closed")

		log.Printf("[server] Client %v connected\n", cl.Id)

		for {
			_, data, err := cl.Conn.Read(ctx)
			if err != nil {
				if websocket.CloseStatus(err) != -1 {
					log.Printf("[info] Client %v disconnected: %v\n", cl.Id, err)
				} else {
					log.Printf("[error] Failed to read connection: %v\n", err)
				}
				return
			}
			log.Printf("[client %v msg] %v\n", cl.Id, string(data))
			go cls.Broadcast(ctx, serverMsg)
		}
	})
}
