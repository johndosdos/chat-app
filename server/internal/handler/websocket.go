package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/johndosdos/chat-app/server/internal/chat"
	"github.com/johndosdos/chat-app/server/internal/database"
	ws "github.com/johndosdos/chat-app/server/internal/websocket"
)

func ServeWs(ctx context.Context, h *ws.Hub, db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("[error] failed to upgrade connection to WebSocket: %v", err)
			return
		}

		// We'll register our new client to the central hub.
		c := ws.NewClient(conn)
		h.Register <- c

		// Load recent chat history to current client.
		go chat.DbLoadChatHistory(ctx, c.Recv, db)

		// Run these goroutines to listen and process messages from other
		// clients.
		go c.WriteMessage()
		go c.ReadMessage()
	}
}
