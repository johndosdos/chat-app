package websocket

import (
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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
		message := <-c.Recv
		err := c.conn.WriteMessage(websocket.TextMessage, message.Content)
		if err != nil {
			log.Printf("[error] failed to write data to client")
			delete(c.Hub.clients, c.userid)
		}
	}
}

func (c *Client) ReadMessage() {
	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("[error] failed to read message from connection: %v", err)
		}

		message := chat.Message{
			Content: p,
			From:    c.userid,
		}
		c.Hub.accept <- message
	}
}
