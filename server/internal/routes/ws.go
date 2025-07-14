package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Printf("Error: failed to accept websocket connection: %v", err)
		return
	}
	defer conn.CloseNow()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data any
	err = wsjson.Read(ctx, conn, data)
	if err != nil {
		log.Printf("Error: failed to read json: %v", err)
		return
	}

	log.Printf("Message: %v", data)

	conn.Close(websocket.StatusNormalClosure, "")
}
