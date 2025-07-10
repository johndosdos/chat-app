package routes

import (
	"fmt"
	"net/http"
)

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Connected to WebSocket!\n`")
}
