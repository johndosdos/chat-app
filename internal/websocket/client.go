package websocket

import (
	"context"
	"encoding/json"
	"log"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	components "github.com/johndosdos/chat-app/server/components/chat"
	"github.com/johndosdos/chat-app/server/internal/chat"
)

type Client struct {
	userid uuid.UUID
	conn   *websocket.Conn
	Hub    *Hub
	Recv   chan chat.Message
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		Recv: make(chan chat.Message),
	}
}

func (c *Client) WriteMessage() {
	for {
		message, ok := <-c.Recv
		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		// Invoke a new writer from the current connection.
		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			log.Printf("[error] %v", err)
			break
		}

		var content templ.Component
		if message.From == c.userid {
			content = components.SenderBubble(message.Content)
		} else {
			content = components.ReceiverBubble(message.Content)
		}
		content.Render(context.Background(), w)

		w.Close()
	}
}

func (c *Client) ReadMessage() {
	defer func() {
		c.Hub.Unregister <- c
		c.conn.Close()
	}()

	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("[error] %v", err)
			}
			break
		}

		// We need to unmarshal the JSON sent from the client side. HTMX's ws-send
		// attribute will also send a HEADERS field along with the client message.
		message := chat.Message{
			From: c.userid,
		}
		err = json.Unmarshal(p, &message)
		if err != nil {
			log.Printf("[error] failed to process payload from client: %v", err)
			break
		}

		c.Hub.accept <- message
	}
}
