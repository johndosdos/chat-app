package clients

import (
	"context"
	"log"
	"sync"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type client struct {
	Id   uuid.UUID
	Conn *websocket.Conn
}

type Clients struct {
	mu   sync.Mutex
	list map[uuid.UUID]client
}

func (cl *client) ReadConn(ctx context.Context) {
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
	}
}

func NewClient(conn *websocket.Conn) client {
	return client{
		Id:   uuid.New(),
		Conn: conn,
	}
}

func NewClients() Clients {
	return Clients{
		list: make(map[uuid.UUID]client),
	}
}

func (cls *Clients) Add(key uuid.UUID, value client) bool {
	cls.mu.Lock()
	defer cls.mu.Unlock()

	_, ok := cls.list[key]
	if !ok {
		cls.list[key] = value
		return true
	}

	return false
}
