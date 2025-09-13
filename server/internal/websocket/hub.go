package websocket

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/johndosdos/chat-app/server/internal/chat"
	"github.com/microcosm-cc/bluemonday"
)

type Hub struct {
	clients    map[uuid.UUID]*Client
	Register   chan *Client
	Unregister chan *Client
	accept     chan chat.Message
	sendToDb   chan chat.Message
	sanitizer  sanitizer
}

type sanitizer interface {
	Sanitize(s string) string
	SanitizeBytes(p []byte) []byte
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case client := <-h.Register:
			client.userid = uuid.New()
			h.clients[client.userid] = client
			client.Hub = h
		case client := <-h.Unregister:
			if _, ok := h.clients[client.userid]; ok {
				delete(h.clients, client.userid)
			}
		case message := <-h.accept:
			// We need to sanitize incoming messages to prevent XSS.
			sanitized := h.sanitizer.SanitizeBytes(message.Content)
			message.Content = sanitized
			h.sendToDb <- message
			for _, client := range h.clients {
				client.Recv <- message
			}
		case <-ctx.Done():
			log.Printf("[error] context cancelled: %v", ctx.Err().Error())
			return
		}
	}
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		accept:     make(chan chat.Message),
		sendToDb:   make(chan chat.Message),
		sanitizer:  bluemonday.StrictPolicy(),
	}
}
