package clients

import (
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
